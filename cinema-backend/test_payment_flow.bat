# Test Script for Payment Flow
# This script tests the complete payment flow

# Step 1: Login to get token
@echo "Step 1: Logging in..."
curl -s -X POST http://localhost:3001/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  > login_response.json

# Extract token (simplified - in real scenario use jq)
@echo "Login completed. Check login_response.json for token."

# Step 2: Create a booking (requires auth token)
@echo "Step 2: To test booking creation, use the token from login_response.json"
@echo "Example:"
@echo "curl -X POST http://localhost:3001/api/bookings/create \"
@echo "  -H \"Authorization: Bearer YOUR_TOKEN\" \"
@echo "  -H \"Content-Type: application/json\" \"
@echo "  -d '{\"showtimeId\":\"YOUR_SHOWTIME_ID\",\"seats\":[{\"seatId\":\"SEAT_ID\",\"seatLabel\":\"A1\",\"price\":10}],\"concessions\":[],\"totalTicketPrice\":10,\"totalConcessionPrice\":0,\"totalAmount\":10}'"

# Step 3: Test mock payment success
@echo ""
@echo "Step 3: Test mock payment success endpoint"
@echo "Visit: http://localhost:3001/api/payments/vnpay/mock-success?bookingId=YOUR_BOOKING_ID"

# Step 4: Check booking status
@echo ""
@echo "Step 4: Check booking status"
@echo "curl http://localhost:3001/api/payments/status/YOUR_BOOKING_ID \"
@echo "  -H \"Authorization: Bearer YOUR_TOKEN\""