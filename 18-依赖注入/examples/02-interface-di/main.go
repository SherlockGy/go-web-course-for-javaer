// 02-interface-di: 基于接口的依赖注入与 Mock 测试
//
// 📌 最佳实践:
//   - 依赖接口而非具体实现
//   - 接口在使用方定义（Go 惯例）
//   - 便于单元测试时 mock
//
// 📌 与 Java 对比:
//   - Java: Mockito.mock(UserRepository.class)
//   - Go: 手动实现接口或使用 gomock
//
// 📌 Go 接口特点:
//   - 隐式实现（无需 implements 关键字）
//   - 接口小而精（1-3 个方法）
//   - "接受接口，返回结构体"
package main

import (
	"errors"
	"fmt"
)

// ==================== 接口定义 ====================

// PaymentGateway 支付网关接口
// 📌 小接口原则：只定义需要的方法
type PaymentGateway interface {
	Charge(userID string, amount float64) (transactionID string, err error)
	Refund(transactionID string) error
}

// NotificationService 通知服务接口
type NotificationService interface {
	SendEmail(to, subject, body string) error
	SendSMS(phone, message string) error
}

// ==================== 真实实现 ====================

// StripeGateway Stripe 支付实现
type StripeGateway struct {
	apiKey string
}

func NewStripeGateway(apiKey string) *StripeGateway {
	return &StripeGateway{apiKey: apiKey}
}

func (g *StripeGateway) Charge(userID string, amount float64) (string, error) {
	// 真实实现会调用 Stripe API
	fmt.Printf("[Stripe] 扣款: 用户=%s, 金额=%.2f\n", userID, amount)
	return "txn_stripe_123", nil
}

func (g *StripeGateway) Refund(transactionID string) error {
	fmt.Printf("[Stripe] 退款: 交易=%s\n", transactionID)
	return nil
}

// TwilioNotification Twilio 通知实现
type TwilioNotification struct {
	accountSID string
	authToken  string
}

func NewTwilioNotification(accountSID, authToken string) *TwilioNotification {
	return &TwilioNotification{accountSID: accountSID, authToken: authToken}
}

func (n *TwilioNotification) SendEmail(to, subject, body string) error {
	fmt.Printf("[Twilio] 发送邮件: to=%s, subject=%s\n", to, subject)
	return nil
}

func (n *TwilioNotification) SendSMS(phone, message string) error {
	fmt.Printf("[Twilio] 发送短信: phone=%s, message=%s\n", phone, message)
	return nil
}

// ==================== Mock 实现（用于测试）====================

// MockPaymentGateway Mock 支付网关
// 📌 测试时使用，可控制返回值
type MockPaymentGateway struct {
	ChargeFunc func(userID string, amount float64) (string, error)
	RefundFunc func(transactionID string) error
	// 记录调用
	ChargeCalled bool
	ChargeArgs   struct {
		UserID string
		Amount float64
	}
}

func (m *MockPaymentGateway) Charge(userID string, amount float64) (string, error) {
	m.ChargeCalled = true
	m.ChargeArgs.UserID = userID
	m.ChargeArgs.Amount = amount

	if m.ChargeFunc != nil {
		return m.ChargeFunc(userID, amount)
	}
	return "mock_txn_123", nil
}

func (m *MockPaymentGateway) Refund(transactionID string) error {
	if m.RefundFunc != nil {
		return m.RefundFunc(transactionID)
	}
	return nil
}

// MockNotificationService Mock 通知服务
type MockNotificationService struct {
	SendEmailFunc func(to, subject, body string) error
	SendSMSFunc   func(phone, message string) error
	EmailsSent    []string
}

func (m *MockNotificationService) SendEmail(to, subject, body string) error {
	m.EmailsSent = append(m.EmailsSent, to)
	if m.SendEmailFunc != nil {
		return m.SendEmailFunc(to, subject, body)
	}
	return nil
}

func (m *MockNotificationService) SendSMS(phone, message string) error {
	if m.SendSMSFunc != nil {
		return m.SendSMSFunc(phone, message)
	}
	return nil
}

// ==================== 业务服务 ====================

// OrderService 订单服务
type OrderService struct {
	payment      PaymentGateway
	notification NotificationService
}

func NewOrderService(payment PaymentGateway, notification NotificationService) *OrderService {
	return &OrderService{
		payment:      payment,
		notification: notification,
	}
}

func (s *OrderService) CreateOrder(userID, email string, amount float64) (string, error) {
	// 1. 扣款
	txnID, err := s.payment.Charge(userID, amount)
	if err != nil {
		return "", fmt.Errorf("支付失败: %w", err)
	}

	// 2. 发送通知
	if err := s.notification.SendEmail(email, "订单确认", fmt.Sprintf("您的订单已创建，交易号: %s", txnID)); err != nil {
		// 记录日志但不影响订单
		fmt.Printf("发送邮件失败: %v\n", err)
	}

	return txnID, nil
}

// ==================== 主函数 ====================

func main() {
	fmt.Println("=== 基于接口的依赖注入 ===\n")

	// 1. 生产环境：使用真实实现
	fmt.Println("--- 生产环境 ---")
	productionPayment := NewStripeGateway("sk_live_xxx")
	productionNotification := NewTwilioNotification("AC123", "token456")
	productionOrderService := NewOrderService(productionPayment, productionNotification)

	txnID, err := productionOrderService.CreateOrder("user_001", "user@example.com", 99.99)
	if err != nil {
		fmt.Printf("创建订单失败: %v\n", err)
		return
	}
	fmt.Printf("订单创建成功: %s\n", txnID)

	// 2. 测试环境：使用 Mock 实现
	fmt.Println("\n--- 测试环境 (Mock) ---")
	mockPayment := &MockPaymentGateway{
		ChargeFunc: func(userID string, amount float64) (string, error) {
			fmt.Printf("[Mock] 模拟扣款: 用户=%s, 金额=%.2f\n", userID, amount)
			return "mock_txn_456", nil
		},
	}
	mockNotification := &MockNotificationService{}

	testOrderService := NewOrderService(mockPayment, mockNotification)
	txnID, _ = testOrderService.CreateOrder("test_user", "test@example.com", 50.00)

	// 验证调用
	fmt.Printf("Mock 支付被调用: %v\n", mockPayment.ChargeCalled)
	fmt.Printf("Mock 支付参数: userID=%s, amount=%.2f\n",
		mockPayment.ChargeArgs.UserID, mockPayment.ChargeArgs.Amount)
	fmt.Printf("Mock 邮件发送列表: %v\n", mockNotification.EmailsSent)

	// 3. 测试失败场景
	fmt.Println("\n--- 测试失败场景 ---")
	failingPayment := &MockPaymentGateway{
		ChargeFunc: func(userID string, amount float64) (string, error) {
			return "", errors.New("余额不足")
		},
	}
	failingOrderService := NewOrderService(failingPayment, mockNotification)
	_, err = failingOrderService.CreateOrder("poor_user", "poor@example.com", 10000.00)
	if err != nil {
		fmt.Printf("预期的失败: %v\n", err)
	}
}
