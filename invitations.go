package main

import (
	"github.com/Take-A-Seat/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func sendInviteEmail(inviteId primitive.ObjectID, userEmail string, codeId string)  {
	var toInfo storage.ToInfo

	toInfo.Email = userEmail

	var subject string
	var textContent string
	var htmlContent string
	subject = "Confirmation reservation"
	textContent = "Take a seat"
	htmlContent = "<p>Your code is:" + codeId

	storage.SendEmail(subject, toInfo, textContent, htmlContent)
}
