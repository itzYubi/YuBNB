go build -o bookings.exe github.com/itzYubi/bookings/cmd/web/.
bookings.exe -dbname=bookings -dbuser=postgres -dbpass=root -cache=false -production=false