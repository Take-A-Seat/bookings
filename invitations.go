package main

func sendConfirmationCreateReservation(userEmail string, firstName string, restaurantName string, hour string, date string) {
	var toInfo ToInfo

	toInfo.Email = userEmail

	var subject string
	var textContent string
	var htmlContent string
	subject = restaurantName + " The reservation is pending"
	textContent = "Take a seat"
	htmlContent = "<p> <strong>" + firstName + "</strong><p> <br> <p>" + "We will get back to you soon with an email confirming your reservation at <strong>" + restaurantName + "</strong>  at  <strong>" + date + " " + hour + "</strong> </p>"

	SendEmail(subject, toInfo, textContent, htmlContent)
}

func sendConfirmationAcceptReservation(userEmail string,firstName string,message string,restaurantName string,code string)  {
	var toInfo ToInfo

	toInfo.Email = userEmail

	var subject string
	var textContent string
	var htmlContent string
	subject = restaurantName + " The reservation was accepted"
	textContent = "Take a seat"
	htmlContent = "<p>Hi <strong>" + firstName + "</strong>, the reservation was accepted<p><br><p>The code is: <strong>"+code+"</strong><p><br><p>Message:" +message+"</p>"

	SendEmail(subject, toInfo, textContent, htmlContent)
}
