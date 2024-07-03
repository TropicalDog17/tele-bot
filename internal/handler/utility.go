package handler

// // Handler for utility buttons like menu, back, close, etc.
// func UtilityHandler(b *tele.Bot, localizer *i18n.Localizer, authRoute *tele.Group, currentStep *string) {
// 	authRoute.Handle(types.BtnBack(localizer), func(c tele.Context) error {
// 		if *currentStep == "confirmOrder" {
// 			*currentStep = ""
// 		}
// 		return c.Send("Back to main menu", types.Menu)
// 	})
// }
