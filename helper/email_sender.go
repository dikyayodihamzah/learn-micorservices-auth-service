package helper

import (
	"crypto/tls"
	"encoding/base64"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/thanhpk/randstr"
	mail "github.com/xhit/go-simple-mail/v2"

	_ "github.com/joho/godotenv/autoload"
)

var (
	host              = os.Getenv("SMTP_HOST")
	port, _           = strconv.Atoi(os.Getenv("SMTP_PORT"))
	emails            = os.Getenv("AUTH_EMAIL")
	password          = os.Getenv("AUTH_PASSWORD")
	sender            = os.Getenv("SENDER_NAME")
	connectTimeout, _ = strconv.Atoi(os.Getenv("CONNECT_TIMEOUT"))
	sendTimeout, _    = strconv.Atoi(os.Getenv("SEND_TIMEOUT"))
	contentText1      = os.Getenv("CONTENT_TEXT1")
	titleReset        = os.Getenv("TITLE_CONTENT_RESET")
	contentTextReset  = os.Getenv("CONTENT_TEXT_RESET")
	buttonTextReset   = os.Getenv("BUTTON_TEXT_RESET")
	regards           = os.Getenv("REGARDS")
	regardsTeam       = os.Getenv("REGARDS_TEAM")
	footerText        = os.Getenv("FOOTER_TEXT")
	poweredText       = os.Getenv("POWERED_TEXT")
	poweredLink       = os.Getenv("POWERED_LINK")
	urlReset          = os.Getenv("URL_RESET_PASSWORD")
	linkPicture       = os.Getenv("LINK_IMAGE_HEAD")
)

func EmailSender(email, token string) error {

	var encodedString = base64.StdEncoding.EncodeToString([]byte(token))
	token = encodedString

	user := emails
	password := password

	url := urlReset + "?email=" + email + "&token=" + token
	subject := "Reset Password #" + strings.ToUpper(randstr.String(10))
	body := `<!doctype html>
	<html lang="en-US">
	
	<head>
	<meta name="viewport" content="width=device-width">
	<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
	<title>Synapsis Email</title>
	<style>
	 @media only screen and (max-width: 620px) {
	 table[class=body] h1 {
	 font-size: 28px !important;
	 margin-bottom: 10px !important;
	 }
	 .image-container {
		display: flex;
		justify-content: center;
	}
	 table[class=body] p,
	 table[class=body] ul,
	 table[class=body] ol,
	 table[class=body] td,
	 table[class=body] span,
	 table[class=body] a {
	 font-size: 16px !important;
	 }
	 table[class=body] .wrapper,
	 table[class=body] .article {
	 padding: 10px !important;
	 }
	 table[class=body] .content {
	 padding: 0 !important;
	 }
	 table[class=body] .container {
	 padding: 0 !important;
	 width: 100% !important;
	 }
	 table[class=body] .main {
	 border-left-width: 0 !important;
	 border-radius: 0 !important;
	 border-right-width: 0 !important;
	 }
	 table[class=body] .btn table {
	 width: 100% !important;
	 }
	 table[class=body] .btn a {
	 width: 100% !important;
	 }
	 table[class=body] .img-responsive {
	 height: auto !important;
	 max-width: 100% !important;
	 width: auto !important;
	 }
	 }
	 @media all {
	 .ExternalClass {
	 width: 100%;
	 }
	 .ExternalClass,
	 .ExternalClass p,
	 .ExternalClass span,
	 .ExternalClass font,
	 .ExternalClass td,
	 .ExternalClass div {
	 line-height: 100%;
	 }
	 .apple-link a {
	 color: inherit !important;
	 font-family: inherit !important;
	 font-size: inherit !important;
	 font-weight: inherit !important;
	 line-height: inherit !important;
	 text-decoration: none !important;
	 }
	 #MessageViewBody a {
	 color: inherit;
	 text-decoration: none;
	 font-size: inherit;
	 font-family: inherit;
	 font-weight: inherit;
	 line-height: inherit;
	 }
	 .btn-primary table td:hover {
	 background-color: #34495e !important;
	 }
	 .btn-primary a:hover {
	 background-color: #34495e !important;
	 border-color: #34495e !important;
	 }
	 }
	</style>
 </head>
 <body class=""
	style="background-color: #f6f6f6; font-family: sans-serif; -webkit-font-smoothing: antialiased; font-size: 14px; line-height: 1.4; margin: 0; padding: 0; -ms-text-size-adjust: 100%; -webkit-text-size-adjust: 100%;">
	<table role="presentation" border="0" cellpadding="0" cellspacing="0" class="body" 
	 style="border-collapse: separate; mso-table-lspace: 0pt; mso-table-rspace: 0pt; width: 100%; background-color: #f6f6f6;">
	 <tr>
		<td style="font-family: sans-serif; font-size: 14px; vertical-align: top;">&nbsp;</td>
		<td class="container"
			 style="font-family: sans-serif; font-size: 14px; vertical-align: top; display: block; Margin: 0 auto; max-width: 580px; padding: 10px; width: 580px;">
			 <!-- START LOGO -->
    		<table role="presentation" class="main" width="200" border="0" cellpadding="0" cellspacing="0" align="center" 
				style="border-collapse: separate; mso-table-lspace: 0pt; mso-table-rspace: 0pt; width: 100%; border-radius: 3px;">
    		    <tbody>
    		        <tr>
    		            <td width="100%" align="center" style="padding-top:10px;">
										<img class="center" src="` + linkPicture + `" width="146" height="66" >
    		            </td>
    		        </tr>
    		    </tbody>
    		</table>
    		<!-- END LOGO -->
			 <div class="content"
				style="box-sizing: border-box; display: block; Margin: 0 auto; max-width: 580px; padding: 10px;">
				<!-- START CENTERED WHITE CONTAINER -->
				<table role="presentation" class="main"
				 style="border-collapse: separate; mso-table-lspace: 0pt; mso-table-rspace: 0pt; width: 100%; background: #ffffff; border-radius: 3px;">
				 <!-- START MAIN CONTENT AREA -->
				 <tr>
					<td class="wrapper"
						 style="font-family: sans-serif; font-size: 14px; vertical-align: top; box-sizing: border-box; padding: 20px;">
						 <table role="presentation" border="0" cellpadding="0" cellspacing="0" align="center"
							style="border-collapse: separate; mso-table-lspace: 0pt; mso-table-rspace: 0pt; width: 100%;">
							<h1 style="color:#1e1e2d; font-weight:500; margin:0;font-size:32px;font-family:'Rubik',sans-serif;text-align: center">
							` + titleReset + `</h1>
							<br>
							<tr>
							 <td align="center" style="font-family: sans-serif; font-size: 14px; vertical-align: top;">
							 <p style="font-family: sans-serif; font-size: 14px; font-weight: normal; margin: 0; text-align: center;">
									 ` + contentText1 + `
								</p>
								<p style="font-family: sans-serif; font-size: 14px; font-weight: normal; margin: 0; Margin-bottom: 15px; text-align: center;text-align: center;">
									 ` + contentTextReset + `
								</p>
								<div style="text-align: center; Margin-bottom: 15px">
								<a href="` + url + `" style="font-family: sans-serif; font-size: 14px; font-weight: normal; margin: 0; Margin-bottom: 15px; text-align: center;">
									 ` + url + `
								</a>
								</div>
								<table role="presentation" border="0" cellpadding="0" cellspacing="0"
									 class="btn btn-primary"
									 style="border-collapse: separate; mso-table-lspace: 0pt; mso-table-rspace: 0pt; width: 100%; box-sizing: border-box;">
									 <tbody>
										<tr>
										 <td align="center"
											style="font-family: sans-serif; font-size: 14px; vertical-align: top; padding-bottom: 15px;">
											<table role="presentation" border="0" cellpadding="0"
												 cellspacing="0"
												 style="border-collapse: separate; mso-table-lspace: 0pt; mso-table-rspace: 0pt; width: auto;">
												 <tbody>
													<tr>
													 <td style="font-family: sans-serif; font-size: 14px; vertical-align: top; horizontal-align: center; background-color: #3498db; border-radius: 5px; text-align: center;">
														<a href="` + url + `" target="_blank"
															 style="display: inline-block; color: #ffffff; background-color: #3498db; border: solid 1px #3498db; border-radius: 5px; box-sizing: border-box; cursor: pointer; text-decoration: none; font-size: 14px; font-weight: bold; margin: 0; padding: 12px 25px; text-transform: capitalize; border-color: #3498db;">` + buttonTextReset + `</a>
													 </td>
													</tr>
												 </tbody>
											</table>
										 </td>
										</tr>
									 </tbody>
								</table>
								<p style="font-family: sans-serif; font-size: 14px; font-weight: normal; margin: 0; Margin-bottom: 15px; text-align: center;">
									 ` + regards + ` <br>` + regardsTeam + `
								</p>
							 </td>
							</tr>
						 </table>
					</td>
				 </tr>
				 <!-- END MAIN CONTENT AREA -->
				</table>
				<!-- START FOOTER -->
				<div class="footer" style="clear: both; Margin-top: 10px; text-align: center; width: 100%;">
				 <table role="presentation" border="0" cellpadding="0" cellspacing="0" align="center"
					style="border-collapse: separate; mso-table-lspace: 0pt; mso-table-rspace: 0pt; width: 100%;">
					<tr>
						 <td class="content-block" align="center"
							style="font-family: sans-serif; vertical-align: top; padding-bottom: 10px; padding-top: 10px; font-size: 12px; color: #999999; text-align: center;">
							<span class="apple-link" style="color: #999999; font-size: 12px; text-align: center;">` + footerText + `</span>.
						 </td>
					</tr>
					<tr>
						 <td class="content-block powered-by" align="center"
							style="font-family: sans-serif; vertical-align: top; padding-bottom: 10px; padding-top: 10px; font-size: 12px; color: #999999; text-align: center;">
							Powered by <a href="` + poweredLink + `"
							 style="color: #999999; font-size: 12px; text-align: center; text-decoration: none;">` + poweredText + `</a>.
						 </td>
					</tr>
				 </table>
				</div>
				<!-- END FOOTER -->
				<!-- END CENTERED WHITE CONTAINER -->
			 </div>
		</td>
		<td style="font-family: sans-serif; font-size: 14px; vertical-align: top;">&nbsp;</td>
	 </tr>
	</table>
 </body>
	
	</html>`

	server := mail.NewSMTPClient()

	server.Host = host
	server.Port = port
	server.Username = user
	server.Password = password
	server.Encryption = mail.EncryptionSTARTTLS

	server.KeepAlive = true

	// Timeout for connect to SMTP Server
	server.ConnectTimeout = time.Duration(connectTimeout) * time.Second

	// Timeout for send the data and wait respond
	server.SendTimeout = time.Duration(sendTimeout) * time.Second

	// Set TLSConfig to provide custom TLS configuration. For example,
	// to skip TLS verification (useful for testing):
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// SMTP client
	smtpClient, err := server.Connect()

	if err != nil {
		log.Println(err)
	}

	emails := mail.NewMSG()
	emails.SetFrom(sender).
		AddTo(email).
		//AddCc("otherto@example.com").
		SetSubject(subject)

	emails.SetBody(mail.TextHTML, body)

	// always check error after send
	if emails.Error != nil {
		log.Println(emails.Error)
	}

	// Call Send and pass the client
	err = emails.Send(smtpClient)
	if err != nil {
		log.Println(err)
		return err
	} else {
		log.Println("Email Sent")
		return nil
	}
}
