package main

func sendConfirmationCreateReservation(userEmail string, firstName string, restaurantName string, startHour string, startDate string, endHour string) {
	var toInfo ToInfo

	toInfo.Email = userEmail

	var subject string
	var textContent string
	var htmlContent string
	subject = restaurantName + " The reservation is pending"
	textContent = "Take a seat"
	htmlContent = "<p> <strong>" + firstName + "</strong><p> <br> <p>" + "We will get back to you soon with an email confirming your reservation at <strong>" + restaurantName + "</strong>  at  <strong>" + startDate + " " + startHour + "</strong> until at " + "</strong>  at <strong>" + endHour + "</p>"

	SendEmail(subject, toInfo, textContent, htmlContent)
}


func sendArrivedClient(userEmail string, firstName string, restaurantName string,code string) {
	var toInfo ToInfo

	toInfo.Email = userEmail

	var subject string
	var textContent string
	var htmlContent string
	subject = restaurantName + " The reservation is pending"
	textContent = "Take a seat"
	htmlContent = "<p> <strong>" + firstName + "</strong><p> <br><p> Your code to access menu is: <strong>"+code+"</strong></p>"

	SendEmail(subject, toInfo, textContent, htmlContent)
}


func sendConfirmationAcceptReservation(userEmail string,firstName string,message string,restaurantName string)  {
	var toInfo ToInfo

	toInfo.Email = userEmail

	var subject string
	var textContent string
	var htmlContent string
	subject = restaurantName + " The reservation was accepted"
	textContent = "Take a seat"
	htmlContent = "<p>Hi <strong>" + firstName + "</strong>, the reservation was accepted<p><br><p>Message from restaurant:" +message+"</p>"

	SendEmail(subject, toInfo, textContent, htmlContent)
}

func sendFinishReservation(userEmail string,firstName string,restaurantName string)  {
	var toInfo ToInfo

	toInfo.Email = userEmail

	var subject string
	var textContent string
	var htmlContent string
	subject = restaurantName + " The reservation is finished"
	textContent = "Take a seat"
	htmlContent = "<p>Hi <strong>" + firstName + "</strong>, <br>Thank you very much!"

	SendEmail(subject, toInfo, textContent, htmlContent)
}

func sendDeclineReservation(userEmail string,firstName string,restaurantName string,message string)  {
	var toInfo ToInfo

	toInfo.Email = userEmail

	var subject string
	var textContent string
	var htmlContent string
	subject = restaurantName + " The reservation is declined"
	textContent = "Take a seat"
	htmlContent = "<p>Hi <strong>" + firstName + "</strong>, <br>Sorry, your reservation is decline"

	SendEmail(subject, toInfo, textContent, htmlContent)
}

