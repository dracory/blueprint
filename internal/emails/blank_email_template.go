package emails

// This file has been commented out as the email functionality has been moved to the base project.
// The new implementation is now using github.com/dracory/base/email package.
// See adapter.go for the new implementation.

// import (
// 	"project/app/links"
// 	"project/config"

// 	"github.com/dromara/carbon/v2"
// 	"github.com/gouniverse/hb"
// )

// // blankEmailTemplate blank HTML template
// func blankEmailTemplate(title string, htmlContent string) string {
// 	copyrightYear := carbon.Now(carbon.UTC).Format("Y")
// 	linkLogin := hb.Hyperlink().
// 		HTML("Login").
// 		Href(links.NewAuthLinks().Login(links.NewUserLinks().Home(map[string]string{}))).
// 		Style("color:white;").
// 		ToHTML()

// 	appName := config.AppName

// 	headerBackGroundColor := "#17A2B8"

// 	template := `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
// <html xmlns="http://www.w3.org/1999/xhtml" style="font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;">
//     <head>
//         <meta name="viewport" content="width=device-width" />
//         <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
//         <title>` + title + `</title>
//         <style type="text/css">
//             img {
//                 max-width: 100%;
//             }
//             body {
//                 -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; width: 100% !important; height: 100%; line-height: 1.6em;
//             }
//             body {
//                 background-color: #f6f6f6 !important;
//             }
//             @media only screen and (max-width: 640px) {
//                 body {
//                     padding: 0 !important;
//                 }
//                 h1 {
//                     font-weight: 800 !important; margin: 20px 0 5px !important;
//                 }
//                 h2 {
//                     font-weight: 800 !important; margin: 20px 0 5px !important;
//                 }
//                 h3 {
//                     font-weight: 800 !important; margin: 20px 0 5px !important;
//                 }
//                 h4 {
//                     font-weight: 800 !important; margin: 20px 0 5px !important;
//                 }
//                 h1 {
//                     font-size: 22px !important;
//                 }
//                 h2 {
//                     font-size: 18px !important;
//                 }
//                 h3 {
//                     font-size: 16px !important;
//                 }
//                 .container {
//                     padding: 0 !important; width: 100% !important;
//                 }
//                 .content {
//                     padding: 0 !important;
//                 }
//                 .content-wrap {
//                     padding: 10px !important;
//                     background: #fff !important;
//                 }
//                 .invoice {
//                     width: 100% !important;
//                 }
//             }
//         </style>
//     </head>

//     <body itemscope itemtype="http://schema.org/EmailMessage" style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; -webkit-font-smoothing: antialiased; -webkit-text-size-adjust: none; width: 100% !important; height: 100%; line-height: 1.6em; background-color: #f6f6f6; margin: 0;" bgcolor="#f6f6f6">
//         <table class="body-wrap" style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; width: 100%; background-color: #f6f6f6; margin: 0;" bgcolor="#f6f6f6">
//             <tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;">
//                 <td style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0;" valign="top"></td>
//                 <td class="container" width="600" style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; display: block !important; max-width: 600px !important; clear: both !important; margin: 0 auto;" valign="top">
//                     <div class="content" style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; max-width: 600px; display: block; margin: 0 auto; padding: 20px;">
//                         <table class="main" width="100%" cellpadding="0" cellspacing="0" style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; border-radius: 3px; background-color: #fff; margin: 0; border: 1px solid #e9e9e9;" bgcolor="#fff">
//                             <tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;">
//                                 <td class="alert alert-warning" style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 16px; vertical-align: top; color: #fff; font-weight: 500; text-align: center; border-radius: 3px 3px 0 0; background-color: ` + headerBackGroundColor + `; margin: 0; padding: 20px;" align="center" bgcolor="#FF9F00" valign="top">
//                                     <table width="100%" border="0">
//                                         <tr>
//                                             <td valign="middle" align="left">
//                                                 <span style="color:white;font-weight: 500;font-size: 24px;">
//                                                     <b>` + appName + `</b>
//                                                 </span>
//                                             </td>
// 											<td valign="middle" align="right" width="100">
// 											` + linkLogin + `
// 											</td>
//                                         </tr>
//                                     </table>
//                                 </td>
//                             </tr>
//                             <tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;">
//                                 <td class="content-wrap" style="background: #fff !important;font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0; padding: 20px;" valign="top">
//                                     <table width="100%" cellpadding="0" cellspacing="0" style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;">
//                                         <tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;">
//                                             <td class="content-block" style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; vertical-align: top; margin: 0; padding: 0 0 20px;" valign="top">
// 										    ` + htmlContent + `
//                                             </td>
//                                         </tr>
//                                     </table>
//                                 </td>
//                             </tr>
//                         </table>
//                         <div class="footer" style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; width: 100%; clear: both; color: #999; margin: 0; padding: 20px;">
//                             <table width="100%" style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;">
//                                 <tr style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 14px; margin: 0;">
//                                     <td class="aligncenter content-block" style="font-family: 'Helvetica Neue',Helvetica,Arial,sans-serif; box-sizing: border-box; font-size: 12px; vertical-align: top; color: #999; text-align: center; margin: 0; padding: 0 0 20px;" align="center" valign="top">
//                                         Copyright &copy; ` + copyrightYear + ` ` + appName + `, All rights reserved.
//                                     </td>
//                                 </tr>
//                             </table>
//                         </div>
//                     </div>
//                 </td>
//             </tr>
//         </table>
//     </body>
// </html>`
// 	return template
// }
