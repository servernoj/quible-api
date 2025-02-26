// Code generated by "jade.go"; DO NOT EDIT.

package emailService

import (
	"bytes"
	"fmt"
	"html"
)

const (
	userActivation__0 = `<!DOCTYPE html><html lang="en" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office"><head><meta charset="utf-8"/><meta http-equiv="x-ua-compatible" content="ie=edge"/><meta name="viewport" content="width=device-width, initial-scale=1"/><meta name="x-apple-disable-message-reformatting"/><style type="text/css">  @import url('https://fonts.googleapis.com/css?family=Merriweather|Open+Sans');

  img {
    border: 0; 
    line-height: 100%; 
    vertical-align: middle;
  }
  .col {
    font-size: 16px; 
    line-height: 25px; 
    vertical-align: top;
  }

  @media screen {
    .col, td, th, div, p {
      font-family: -apple-system,system-ui,BlinkMacSystemFont,"Segoe UI","Roboto","Helvetica Neue",Arial,sans-serif;
    }
    .sans-serif {
      font-family: 'Open Sans', Arial, sans-serif;
    }
    .serif {
      font-family: 'Merriweather', Georgia, serif;
    }
    img {
      max-width: 100%;
    }
  }

  @media (max-width: 632px) {
    .container {
      width: 100%!important;
    }
  }

  @media (max-width: 480px) {
    .col {
      display: inline-block!important;
      line-height: 23px;
      width: 100%!important;
    }
    .col-sm-1 {
      max-width: 25%;
    }
    .col-sm-2 {
      max-width: 50%;
    }
    .col-sm-3 {
      max-width: 75%;
    }
    .col-sm-third {
      max-width: 33.33333%;
    }
    .col-sm-push-1 {
      margin-left: 25%;
    }
    .col-sm-push-2 {
      margin-left: 50%;
    }
    .col-sm-push-3 {
      margin-left: 75%;
    }
    .col-sm-push-third {
      margin-left: 33.33333%;
    }
    .full-width-sm {
      display: table!important; 
      width: 100%!important;
    }
    .stack-sm-first {
      display: table-header-group!important;
    }
    .stack-sm-last {
      display: table-footer-group!important;
    }
    .stack-sm-top {
      display: table-caption!important; 
      max-width: 100%; 
      padding-left: 0!important;
    }
    .toggle-content {
      max-height: 0;
      overflow: auto;
      transition: max-height .4s linear;
      -webkit-transition: max-height .4s linear;
    }
    .toggle-trigger:hover + .toggle-content,
    .toggle-content:hover {
      max-height: 999px!important;
    }
    .show-sm {
      display: inherit!important;
      font-size: inherit!important;
      line-height: inherit!important;
      max-height: none!important;
    }
    .hide-sm {
      display: none!important;
    }
    .align-sm-center {
      display: table!important;
      float: none;
      margin-left: auto!important;
      margin-right: auto!important;
    }
    .align-sm-left {
      float: left;
    }
    .align-sm-right {
      float: right;
    }
    .text-sm-center {
      text-align: center!important;
    }
    .text-sm-left {
      text-align: left!important;
    }
    .text-sm-right {
      text-align: right!important;
    }
    .borderless-sm {
      border: none!important;
    }
    .nav-sm-vertical .nav-item {
      display: block;
    }
    .nav-sm-vertical .nav-item a {
      display: inline-block; 
      padding: 4px 0!important;
    }
    .spacer {
      height: 0;
    }
    .p-sm-0 {
      padding: 0!important;
    }
    .p-sm-8 {
      padding: 8px!important;
    }
    .p-sm-16 {
      padding: 16px!important;
    }
    .p-sm-24 {
      padding: 24px!important;
    }
    .pt-sm-0 {
      padding-top: 0!important;
    }
    .pt-sm-8 {
      padding-top: 8px!important;
    }
    .pt-sm-16 {
      padding-top: 16px!important;
    }
    .pt-sm-24 {
      padding-top: 24px!important;
    }
    .pr-sm-0 {
      padding-right: 0!important;
    }
    .pr-sm-8 {
      padding-right: 8px!important;
    }
    .pr-sm-16 {
      padding-right: 16px!important;
    }
    .pr-sm-24 {
      padding-right: 24px!important;
    }
    .pb-sm-0 {
      padding-bottom: 0!important;
    }
    .pb-sm-8 {
      padding-bottom: 8px!important;
    }
    .pb-sm-16 {
      padding-bottom: 16px!important;
    }
    .pb-sm-24 {
      padding-bottom: 24px!important;
    }
    .pl-sm-0 {
      padding-left: 0!important;
    }
    .pl-sm-8 {
      padding-left: 8px!important;
    }
    .pl-sm-16 {
      padding-left: 16px!important;
    }
    .pl-sm-24 {
      padding-left: 24px!important;
    }
    .px-sm-0 {
      padding-right: 0!important; 
      padding-left: 0!important;
    }
    .px-sm-8 {
      padding-right: 8px!important; 
      padding-left: 8px!important;
    }
    .px-sm-16 {
      padding-right: 16px!important; 
      padding-left: 16px!important;
    }
    .px-sm-24 {
      padding-right: 24px!important; 
      padding-left: 24px!important;
    }
    .py-sm-0 {
      padding-top: 0!important; 
      padding-bottom: 0!important;
    }
    .py-sm-8 {
      padding-top: 8px!important; 
      padding-bottom: 8px!important;
    }
    .py-sm-16 {
      padding-top: 16px!important; 
      padding-bottom: 16px!important;
    }
    .py-sm-24 {
      padding-top: 24px!important; 
      padding-bottom: 24px!important;
    }
  }</style></head><body style="margin:0;padding:0;width:100%;word-break:break-word;-webkit-font-smoothing:antialiased;"><div style="display:none;font-size:0;line-height:0;"></div>`
	userActivation__1  = `</body></html>`
	userActivation__2  = `<table lang="en" bgcolor="`
	userActivation__3  = `" cellpadding="16" cellspacing="0" role="presentation" width="100%"><tr><td align="center">`
	userActivation__4  = `</td></tr></table>`
	userActivation__5  = `<table class="container" bgcolor="`
	userActivation__6  = `" cellpadding="0" cellspacing="0" role="presentation" width="600"><tr><td align="left">`
	userActivation__11 = `<h3>Hi `
	userActivation__12 = `, </h3><p>Thank you for signing up for Quible. Click on the link below to verify your email:</p>`
	userActivation__13 = `<p>This link will expire in 24 hours. Feel free to repeat registration to get a new link should this one expire. </p><p>If you did not sign up for a Quible account, you can safely ignore this email. </p><p>Best,</p><p>The Quible Team </p>`
	userActivation__14 = `<a href="`
	userActivation__15 = `" style="`
	userActivation__16 = `">`
	userActivation__17 = `</a>`
)

func UserActivation(name string, link string, buffer *bytes.Buffer) {

	buffer.WriteString(userActivation__0)

	{
		var (
			bg = "#FFF"
		)
		var block []byte
		{
			buffer := new(bytes.Buffer)
			{
				var (
					bg = "#FFF"
				)
				var block []byte
				{
					buffer := new(bytes.Buffer)
					{
						var (
							bg = "#FFF"
						)
						var block []byte
						{
							buffer := new(bytes.Buffer)
							buffer.WriteString(userActivation__11)
							buffer.WriteString(html.EscapeString(fmt.Sprintf("%v", name)))
							buffer.WriteString(userActivation__12)

							{
								var (
									url = link
									fg  = "rgb(17, 85, 204)"
								)
								var block []byte
								{
									buffer := new(bytes.Buffer)
									buffer.WriteString(html.EscapeString(fmt.Sprintf("%v", link)))
									block = buffer.Bytes()
								}

								buffer.WriteString(userActivation__14)
								buffer.WriteString(html.EscapeString(fmt.Sprintf("%v", url)))
								buffer.WriteString(userActivation__15)
								buffer.WriteString(html.EscapeString(fmt.Sprintf("%v", "color: "+fg+"; display: inline-block; line-height: 100%; text-decoration: none;")))
								buffer.WriteString(userActivation__16)
								buffer.Write(block)
								buffer.WriteString(userActivation__17)
							}

							buffer.WriteString(userActivation__13)

							block = buffer.Bytes()
						}

						buffer.WriteString(userActivation__5)
						buffer.WriteString(html.EscapeString(fmt.Sprintf("%v", bg)))
						buffer.WriteString(userActivation__6)

						buffer.Write(block)
						buffer.WriteString(userActivation__4)

					}

					block = buffer.Bytes()
				}

				buffer.WriteString(userActivation__5)
				buffer.WriteString(html.EscapeString(fmt.Sprintf("%v", bg)))
				buffer.WriteString(userActivation__6)

				buffer.Write(block)
				buffer.WriteString(userActivation__4)

			}

			block = buffer.Bytes()
		}

		buffer.WriteString(userActivation__2)
		buffer.WriteString(html.EscapeString(fmt.Sprintf("%v", bg)))
		buffer.WriteString(userActivation__3)

		buffer.Write(block)
		buffer.WriteString(userActivation__4)

	}

	buffer.WriteString(userActivation__1)

}
