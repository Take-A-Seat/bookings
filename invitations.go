package main

func sendConfirmationCreateReservation(userEmail string, firstName string, restaurantName string) {
	var toInfo ToInfo

	toInfo.Email = userEmail

	var subject string
	var textContent string
	var htmlContent string
	subject = restaurantName+" The reservation is pending"
	textContent = "Take a seat"
	htmlContent = "<p> <strong>" + firstName + "</strong><p> <br> <p>" + "We will get back to you soon with an email confirming your reservation at <strong>" + restaurantName + "</strong> </p>"

	SendEmail(subject, toInfo, textContent, htmlContent)
}
