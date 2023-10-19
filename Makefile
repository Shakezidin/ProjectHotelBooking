run:
	go run main.go

test:
	golint routes/AdminRoutes.go
	golint routes/ownerRoutes.go
	golint routes/UserRoutes.go
	golint models/Admin.go
	golint models/AdminRevenue.go
	golint models/Banner.go
	golint models/BookingModel.go
	golint models/Contact.go
	golint models/Coupen.go
	golint models/HashPassword.go
	golint models/HotelModel.go
	golint models/Login.go
	golint models/OwnerModel.go
	golint models/Report.go
	golint models/RoomModel.go
	golint models/Searching.go
	golint models/UserModel.go
	golint middleware/AdminMiddleWare.go
	golint middleware/OwnerMiddleWare.go
	golint middleware/UserMiddleWare.go
	golint initializer/db.go
	golint initializer/env.go
	golint initializer/Redis.go
	golint helper/Available.go
	golint helper/RandString.go
	golint controllers/Admin/AdminHome.go
	golint controllers/Admin/AdminLogin.go
	golint controllers/Admin/BannerManagement.go
	golint controllers/Admin/Cancellationmanagement.go
	golint controllers/Admin/CatagoriesManagement.go
	golint controllers/Admin/CouponManagement.go
	golint controllers/Admin/FecilitiesManagement.go
	golint controllers/Admin/HotelManage.go
	golint controllers/Admin/MessageManagement.go
	golint controllers/Admin/Ownermanage.go
	golint controllers/Admin/ReportManagement.go
	golint controllers/Admin/RoomCatagoriesmanagemnt.go
	golint controllers/Admin/RoomFecilitiesManage.go
	golint controllers/Admin/RoomManage.go
	golint controllers/Admin/Usermanagement.go
	golint controllers/hotelowner/BannerManagement.go
	golint controllers/hotelowner/Dashbord.go
	golint controllers/hotelowner/HotelManagment.go
	golint controllers/hotelowner/OwnerLogin.go
	golint controllers/hotelowner/RoomManagement.go
	golint controllers/Otp/Otplog.go
	golint controllers/User/Booking.go
	golint controllers/User/Cancellation.go
	golint controllers/User/Contact.go
	golint controllers/User/Coupon.go
	golint controllers/User/HotelManagement.go
	golint controllers/User/Razorpay.go
	golint controllers/User/RoomFilter.go
	golint controllers/User/RoomManagement.go
	golint controllers/User/UserHome.go
	golint controllers/User/userLogin.go
	golint controllers/User/UserProfile.go
	golint controllers/User/userSignup.go
	golint Auth/jwt.go



