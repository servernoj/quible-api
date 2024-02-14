// Code generated by "jade.go"; DO NOT EDIT.

package emailService

import (
	"bytes"
	"fmt"
	"html"
)

const (
	demo__0 = `<!DOCTYPE html><html lang="en" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office"><head><meta charset="utf-8"/><meta http-equiv="x-ua-compatible" content="ie=edge"/><meta name="viewport" content="width=device-width, initial-scale=1"/><meta name="x-apple-disable-message-reformatting"/><style type="text/css">  @import url('https://fonts.googleapis.com/css?family=Merriweather|Open+Sans');

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
	demo__1  = `</body></html>`
	demo__2  = `<table lang="en" bgcolor="`
	demo__3  = `" cellpadding="16" cellspacing="0" role="presentation" width="100%"><tr><td align="center">`
	demo__4  = `</td></tr></table>`
	demo__5  = `<table class="container" bgcolor="`
	demo__6  = `" cellpadding="0" cellspacing="0" role="presentation" width="600"><tr><td align="left">`
	demo__11 = `<img src="`
	demo__12 = `" width="100%"/><h1>Hi `
	demo__13 = `,</h1><p>A long paragraph... Lorem ipsum dolor sit amet consectetur adipisicing elit. Obcaecati reprehenderit sed voluptatum nulla ipsa. Necessitatibus, fuga. Animi quod ab dolores corrupti similique incidunt? Aperiam deserunt non cumque veritatis excepturi voluptatibus!</p>`
	demo__14 = `<table cellpadding="0" cellspacing="0" role="presentation" width="100%"><tr><td class="spacer" height="`
	demo__15 = `">`
	demo__16 = `<span>&nbsp;</span></td></tr></table>`
	demo__20 = `<h2 style="color: rgb(6,144,250);">New Tools</h2>`
	demo__21 = `<ul style="list-style: none; padding: 0;">`
	demo__22 = `</ul>`
	demo__23 = `<li>`
	demo__24 = `</li>`
	demo__25 = `<table cellpadding="0" cellspacing="0" role="presentation" width="100%"><tr><td style="`
	demo__26 = `"><table cellpadding="0" cellspacing="0" role="presentation" width="100%"><tr class="full-width-sm">`
	demo__27 = `</tr></table></td></tr></table>`
	demo__28 = `<td width="138" align="left" style="vertical-align: top;"><img src="`
	demo__29 = `" width="75%" style="max-width: 100px"/></td>`
	demo__30 = `<td class="`
	demo__31 = `" width="414">`
	demo__32 = `</td>`
	demo__33 = `<h3 style="color: rgb(6,144,250); padding: 0; margin: 0;">`
	demo__34 = `</h3><p>`
	demo__35 = `</p>`
)

func Demo(imageURL string, name string, buffer *bytes.Buffer) {

	buffer.WriteString(demo__0)

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
							buffer.WriteString(demo__11)
							buffer.WriteString(html.EscapeString(fmt.Sprintf("%v", imageURL)))
							buffer.WriteString(demo__12)
							buffer.WriteString(html.EscapeString(fmt.Sprintf("%v", name)))
							buffer.WriteString(demo__13)

							block = buffer.Bytes()
						}

						buffer.WriteString(demo__5)
						buffer.WriteString(html.EscapeString(fmt.Sprintf("%v", bg)))
						buffer.WriteString(demo__6)

						buffer.Write(block)
						buffer.WriteString(demo__4)

					}

					{
						var (
							height = 16
						)
						var block []byte
						buffer.WriteString(demo__14)
						buffer.WriteString(html.EscapeString(fmt.Sprintf("%v", height)))
						buffer.WriteString(demo__15)
						buffer.Write(block)
						buffer.WriteString(demo__16)

					}

					{
						var (
							bg = "#FFF"
						)
						var block []byte
						{
							buffer := new(bytes.Buffer)
							buffer.WriteString(demo__20)

							data := []struct {
								title       string
								description string
								imageUrl    string
							}{
								{
									"SuperDuperDB",
									"Open-source framework for integrating AI with databases",
									"https://img.stackshare.io/service/145319/default_757ebbcea223c42d3133167ba630eac290704cad.jpg",
								},
								{
									"WarpBuild",
									"x86-64 and arm GitHub Action runners for 30% faster builds",
									"https://img.stackshare.io/service/145320/default_50c6f1b1e1e030ad80920774a4bb9030d5ba7a96.png",
								},
							}
							buffer.WriteString(demo__21)
							for _, item := range data {
								buffer.WriteString(demo__23)
								{
									var (
										padding = "0 0"
									)
									var block []byte
									{
										buffer := new(bytes.Buffer)
										buffer.WriteString(demo__28)
										buffer.WriteString(html.EscapeString(fmt.Sprintf("%v", item.imageUrl)))
										buffer.WriteString(demo__29)

										{
											var (
												isCol = true
											)
											var block []byte
											{
												buffer := new(bytes.Buffer)
												buffer.WriteString(demo__33)
												buffer.WriteString(html.EscapeString(fmt.Sprintf("%v", item.title)))
												buffer.WriteString(demo__34)
												buffer.WriteString(html.EscapeString(fmt.Sprintf("%v", item.description)))
												buffer.WriteString(demo__35)
												block = buffer.Bytes()
											}

											buffer.WriteString(demo__30)
											buffer.WriteString(html.EscapeString(fmt.Sprintf("%v", ternary(isCol, "col", ""))))
											buffer.WriteString(demo__31)
											buffer.Write(block)
											buffer.WriteString(demo__32)
										}

										block = buffer.Bytes()
									}

									buffer.WriteString(demo__25)
									buffer.WriteString(html.EscapeString(fmt.Sprintf("%v", "padding: "+padding+";")))
									buffer.WriteString(demo__26)

									buffer.Write(block)
									buffer.WriteString(demo__27)

								}

								{
									var (
										height = 16
									)
									var block []byte
									buffer.WriteString(demo__14)
									buffer.WriteString(html.EscapeString(fmt.Sprintf("%v", height)))
									buffer.WriteString(demo__15)
									buffer.Write(block)
									buffer.WriteString(demo__16)

								}

								buffer.WriteString(demo__24)
							}
							buffer.WriteString(demo__22)
							block = buffer.Bytes()
						}

						buffer.WriteString(demo__5)
						buffer.WriteString(html.EscapeString(fmt.Sprintf("%v", bg)))
						buffer.WriteString(demo__6)

						buffer.Write(block)
						buffer.WriteString(demo__4)

					}

					block = buffer.Bytes()
				}

				buffer.WriteString(demo__5)
				buffer.WriteString(html.EscapeString(fmt.Sprintf("%v", bg)))
				buffer.WriteString(demo__6)

				buffer.Write(block)
				buffer.WriteString(demo__4)

			}

			block = buffer.Bytes()
		}

		buffer.WriteString(demo__2)
		buffer.WriteString(html.EscapeString(fmt.Sprintf("%v", bg)))
		buffer.WriteString(demo__3)

		buffer.Write(block)
		buffer.WriteString(demo__4)

	}

	buffer.WriteString(demo__1)

}
