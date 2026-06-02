# Luồng Thanh Toán Thành Công - Kiểm Tra

## Tổng quan luồng

```
User chọn ghế 
  → Tạo Booking (Status: PENDING, PaymentStatus: PENDING)
  → Tạo Payment Record (Status: PENDING)
  → Redirect sang VNPay
  → VNPay xử lý thanh toán
  → VNPay callback về /api/payments/vnpay/return
  → Backend verify chữ ký
  → Nếu thành công:
    1. Cập nhật Booking status = CONFIRMED
    2. Cập nhật Payment status = PAID
    3. Tạo Ticket record
    4. Tạo QR Code
    5. Gửi email xác nhận
    6. Redirect về /payment/callback?status=success
  → Frontend hiển thị thông báo thành công
  → Auto redirect sang /customer/tickets sau 5 giây
```

## Các bước kiểm tra

### 1. Tạo Booking
- **Endpoint**: POST /api/bookings/create
- **Headers**: Authorization: Bearer {token}
- **Body**:
```json
{
  "showtimeId": "uuid",
  "seats": [
    {"seatId": "uuid", "seatLabel": "A1", "price": 10.00}
  ],
  "concessions": [],
  "totalTicketPrice": 10.00,
  "totalConcessionPrice": 0,
  "totalAmount": 10.00
}
```
- **Expected Response**:
```json
{
  "success": true,
  "data": {
    "id": "booking-uuid",
    "bookingCode": "BK...",
    "status": "pending",
    "paymentStatus": "pending",
    ...
  }
}
```

### 2. Tạo Payment URL
- **Endpoint**: POST /api/payments/vnpay/create
- **Headers**: Authorization: Bearer {token}
- **Body**:
```json
{
  "bookingId": "booking-uuid",
  "amount": 10.00,
  "orderInfo": "Payment for movie..."
}
```
- **Expected Response**:
```json
{
  "success": true,
  "data": {
    "paymentUrl": "https://sandbox.vnpayment.vn/...",
    "orderId": "order-uuid",
    "paymentId": "payment-uuid"
  }
}
```

### 3. Mock Payment Success (Test)
- **URL**: http://localhost:3001/api/payments/vnpay/mock-success?bookingId={bookingId}&orderId={orderId}
- **Expected**: Redirect về http://localhost:3000/payment/callback?status=success&bookingId={bookingId}

### 4. Kiểm tra Booking sau thanh toán
- **Endpoint**: GET /api/payments/status/{bookingId}
- **Headers**: Authorization: Bearer {token}
- **Expected Response**:
```json
{
  "success": true,
  "data": {
    "bookingId": "booking-uuid",
    "status": "confirmed",
    "paymentStatus": "paid",
    "amount": 10.00,
    "bookingCode": "BK...",
    "qrCode": "QR-..."
  }
}
```

### 5. Kiểm tra Ticket được tạo
- **Endpoint**: GET /api/tickets/my-tickets
- **Headers**: Authorization: Bearer {token}
- **Expected**: Ticket mới với status = "paid" và qrCode

### 6. Kiểm tra Email
- Nếu SMTP được cấu hình, email xác nhận sẽ được gửi
- Nếu không, log sẽ hiển thị: "Email service not configured. Skipping email send."

## Cấu hình Email (Tùy chọn)

Thêm vào file `.env`:
```
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
FROM_EMAIL=your-email@gmail.com
```

## Xử lý lỗi

### Payment Failed
- Booking status = CANCELLED
- Payment status = FAILED
- User có thể thử lại

### Payment Expired
- Booking status = EXPIRED
- Payment status = FAILED
- Ghế được giải phóng
- User cần đặt lại

## Luồng Frontend

### Booking Flow (Step 5 - Payment)
1. User click "Pay with VNPay"
2. Frontend gọi POST /api/payments/vnpay/create
3. Nhận paymentUrl
4. Lưu bookingId vào sessionStorage
5. Redirect sang paymentUrl

### Payment Callback Page
1. Nhận query params từ VNPay
2. Nếu status=success:
   - Hiển thị thông báo thành công
   - Gọi API kiểm tra booking details
   - Hiển thị booking code, movie info, QR code
   - Countdown 5 giây
   - Auto redirect sang /customer/tickets
3. Nếu status=failed:
   - Hiển thị thông báo lỗi
   - Cho phép thử lại

## My Tickets Page
- Hiển thị tất cả bookings
- Status colors:
  - PENDING: Yellow
  - CONFIRMED: Green
  - CANCELLED: Red
  - EXPIRED: Gray
  - COMPLETED: Blue
- CONFIRMED bookings có nút "Show QR"
- PENDING bookings có nút "Pay Now"

## API Endpoints

### Public
- GET /api/payments/vnpay/return - VNPay callback
- GET /api/payments/vnpay/mock - Mock payment page
- GET /api/payments/vnpay/mock-success - Mock success
- GET /api/payments/vnpay/mock-fail - Mock fail

### Protected (Customer)
- POST /api/payments/vnpay/create - Tạo payment URL
- GET /api/payments/status/:bookingId - Kiểm tra status
- POST /api/bookings/create - Tạo booking
- GET /api/my-bookings - Lấy bookings của user
- GET /api/tickets/my-tickets - Lấy tickets của user

### Protected (Admin)
- GET /api/bookings - Lấy tất cả bookings
- POST /api/bookings/:id/confirm - Confirm booking
- POST /api/bookings/:id/cancel - Cancel booking
- POST /api/bookings/:id/refund - Refund booking