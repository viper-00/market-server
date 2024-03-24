package mail

import "market/global"

func UserRegisterTemplate(to string, loginUrl string) []byte {
	msg := []byte(
		"From: " + global.MARKET_CONFIG.Smtp.Username + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: Verify your email\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"utf-8\"\r\n" +
			"\r\n" +
			"<html><body>" +
			`<center>
			<table
			  border="0"
			  cellpadding="0"
			  cellspacing="0"
			  height="100%"
			  width="100%"
			  id="bodyTable"
			  style="background-color: rgb(244, 244, 244)"
			>
			  <tbody>
				<tr>
				  <td class="bodyCell" align="center" valign="top">
					<table id="root" border="0" cellpadding="0" cellspacing="0" width="100%">
					  <tbody data-block-id="15" class="mceWrapper">
						<tr>
						  <td align="center" valign="top" class="mceWrapperOuter">
							<table
							  border="0"
							  cellpadding="0"
							  cellspacing="0"
							  width="100%"
							  style="max-width: 660px"
							  role="presentation"
							>
							  <tbody>
								<tr>
								  <td
									style="
									  background-color: #ffffff;
									  background-position: center;
									  background-repeat: no-repeat;
									  background-size: cover;
									"
									class="mceWrapperInner"
									valign="top"
								  >
									<table
									  align="center"
									  border="0"
									  cellpadding="0"
									  cellspacing="0"
									  width="100%"
									  role="presentation"
									  data-block-id="14"
									>
									  <tbody>
										<tr class="mceRow">
										  <td
											style="
											  background-position: center;
											  background-repeat: no-repeat;
											  background-size: cover;
											"
											valign="top"
										  >
											<table border="0" cellpadding="0" cellspacing="0" width="100%" role="presentation">
											  <tbody>
												<tr>
												  <td
													style="padding-top: 0; padding-bottom: 0"
													class="mceColumn"
													data-block-id="-4"
													valign="top"
													colspan="12"
													width="100%"
												  >
													<table
													  border="0"
													  cellpadding="0"
													  cellspacing="0"
													  width="100%"
													  role="presentation"
													>
													  <tbody>
														<tr>
														  <td
															style="
															  padding-top: 12px;
															  padding-bottom: 12px;
															  padding-right: 48px;
															  padding-left: 48px;
															"
															class="mceBlockContainer"
															align="center"
															valign="top"
														  >
															<img
															  data-block-id="3"
															  width="150"
															  height="auto"
															  style="
																width: 150px;
																height: auto;
																max-width: 150px !important;
																display: block;
															  "
															  alt="Logo"
															  src="https://mcusercontent.com/beb0655d41009e40b7fb5a34d/images/86364ab7-2313-6890-abee-3d6095ca61ef.png"
															  class="mceLogo"
															/>
														  </td>
														</tr>
														<tr>
														  <td
															style="
															  padding-top: 12px;
															  padding-bottom: 12px;
															  padding-right: 24px;
															  padding-left: 24px;
															"
															class="mceBlockContainer"
															valign="top"
														  >
															<div
															  data-block-id="5"
															  class="mceText"
															  id="dataBlockId-5"
															  style="width: 100%"
															>
															  <h1 class="last-child" style="text-align: center">
																<span style="font-size: 30px">Verify you email address</span>
															  </h1>
															</div>
														  </td>
														</tr>
														<tr>
														  <td
															style="
															  padding-top: 12px;
															  padding-bottom: 12px;
															  padding-right: 100px;
															  padding-left: 100px;
															"
															class="mceBlockContainer"
															valign="top"
														  >
															<div
															  data-block-id="6"
															  class="mceText"
															  id="dataBlockId-6"
															  style="width: 100%"
															>
															  <p style="text-align: left">Hi guy,</p>
															  <p style="text-align: left">
																For your security, please verify your email address so we can
																finish setting up you Predictmarket account.
															  </p>
															  <p style="text-align: left" class="last-child">
																If you didn’t register for an account with this email address,
																please ignore this email.
															  </p>
															</div>
														  </td>
														</tr>
														<tr>
														  <td
															style="
															  padding-top: 30px;
															  padding-bottom: 12px;
															  padding-right: 0;
															  padding-left: 0;
															"
															class="mceBlockContainer"
															align="center"
															valign="top"
														  >
															<table
															  align="center"
															  border="0"
															  cellpadding="0"
															  cellspacing="0"
															  role="presentation"
															  data-block-id="7"
															>
															  <tbody>
																<tr class="mceStandardButton">
																  <td
																	style="
																	  background-color: #000000;
																	  border-radius: 4px;
																	  text-align: center;
																	"
																	class="mceButton"
																	valign="top"
																  >
																	<a
																	  href="` + loginUrl + `"
																	  target="_blank"
																	  style="
																		background-color: #000000;
																		border-radius: 4px;
																		border: 2px solid #000000;
																		color: #ffffff;
																		display: block;
																		font-family: 'Helvetica Neue', Helvetica, Arial, Verdana,
																		  sans-serif;
																		font-size: 16px;
																		font-weight: normal;
																		font-style: normal;
																		padding: 16px 28px;
																		text-decoration: none;
																		min-width: 30px;
																		text-align: center;
																		direction: ltr;
																		letter-spacing: 0px;
																	  "
																	  >VERIFY EMAIL</a
																	>
																  </td>
																</tr>
																<tr>
																  <!--[if mso]>
																	<td align="center">
																	  <v:roundrect
																		xmlns:v="urn:schemas-microsoft-com:vml"
																		xmlns:w="urn:schemas-microsoft-com:office:word"
																		href=""
																		style="
																		  v-text-anchor: middle;
																		  width: 165.78px;
																		  height: 55px;
																		"
																		arcsize="2%"
																		strokecolor="#000000"
																		strokeweight="2px"
																		fillcolor="#000000"
																	  >
																		<v:stroke dashstyle="solid" />
																		<w:anchorlock />
																		<center
																		  style="
																			color: #ffffff;
																			display: block;
																			font-family: 'Helvetica Neue', Helvetica, Arial,
																			  Verdana, sans-serif;
																			font-size: 16;
																			font-style: normal;
																			font-weight: normal;
																			letter-spacing: 0px;
																			text-decoration: none;
																			text-align: center;
																			direction: ltr;
																		  "
																		>
																		  VERIFY EMAIL
																		</center>
																	  </v:roundrect>
																	</td>
																  <![endif]-->
																</tr>
															  </tbody>
															</table>
														  </td>
														</tr>
														<tr>
														  <td
															style="
															  padding-top: 8px;
															  padding-bottom: 8px;
															  padding-right: 8px;
															  padding-left: 8px;
															"
															class="mceLayoutContainer"
															valign="top"
														  >
															<table
															  align="center"
															  border="0"
															  cellpadding="0"
															  cellspacing="0"
															  width="100%"
															  role="presentation"
															  data-block-id="13"
															  id="section_37951e71c1cf228f036a7cc5ce1086cc"
															  class="mceFooterSection"
															>
															  <tbody>
																<tr class="mceRow">
																  <td
																	style="
																	  background-position: center;
																	  background-repeat: no-repeat;
																	  background-size: cover;
																	"
																	valign="top"
																  >
																	<table
																	  border="0"
																	  cellpadding="0"
																	  cellspacing="12"
																	  width="100%"
																	  role="presentation"
																	>
																	  <tbody>
																		<tr>
																		  <td
																			style="
																			  padding-top: 0;
																			  padding-bottom: 0;
																			  margin-bottom: 12px;
																			"
																			class="mceColumn"
																			data-block-id="-3"
																			valign="top"
																			colspan="12"
																			width="100%"
																		  >
																			<table
																			  border="0"
																			  cellpadding="0"
																			  cellspacing="0"
																			  width="100%"
																			  role="presentation"
																			>
																			  <tbody>
																				<tr>
																				  <td
																					style="
																					  padding-top: 12px;
																					  padding-bottom: 12px;
																					  padding-right: 16px;
																					  padding-left: 16px;
																					"
																					class="mceBlockContainer"
																					align="center"
																					valign="top"
																				  >
																					<div
																					  data-block-id="11"
																					  class="mceText"
																					  id="dataBlockId-11"
																					  style="display: inline-block; width: 100%"
																					>
																					  <p class="last-child">
																						<em
																						  ><span style="font-size: 12px"
																							>@2024 PredictMarket Inc.</span
																						  ></em
																						>
																					  </p>
																					</div>
																				  </td>
																				</tr>
																			  </tbody>
																			</table>
																		  </td>
																		</tr>
																	  </tbody>
																	</table>
																  </td>
																</tr>
															  </tbody>
															</table>
														  </td>
														</tr>
													  </tbody>
													</table>
												  </td>
												</tr>
											  </tbody>
											</table>
										  </td>
										</tr>
									  </tbody>
									</table>
								  </td>
								</tr>
							  </tbody>
							</table>
						  </td>
						</tr>
					  </tbody>
					</table>
				  </td>
				</tr>
			  </tbody>
			</table>
		  </center>
		  ` +
			"</body></html>\r\n")
	return msg
}

func UserLoginTemplate(username, to, code string) []byte {
	msg := []byte(
		"From: " + global.MARKET_CONFIG.Smtp.Username + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: Login to Predictmarket\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"utf-8\"\r\n" +
			"\r\n" +
			"<html><body>" +
			`<center>
			<table
			  border="0"
			  cellpadding="0"
			  cellspacing="0"
			  height="100%"
			  width="100%"
			  id="bodyTable"
			  style="background-color: rgb(244, 244, 244)"
			>
			  <tbody>
				<tr>
				  <td class="bodyCell" align="center" valign="top">
					<table id="root" border="0" cellpadding="0" cellspacing="0" width="100%">
					  <tbody data-block-id="15" class="mceWrapper">
						<tr>
						  <td align="center" valign="top" class="mceWrapperOuter">
							<table
							  border="0"
							  cellpadding="0"
							  cellspacing="0"
							  width="100%"
							  style="max-width: 660px"
							  role="presentation"
							>
							  <tbody>
								<tr>
								  <td
									style="
									  background-color: #ffffff;
									  background-position: center;
									  background-repeat: no-repeat;
									  background-size: cover;
									"
									class="mceWrapperInner"
									valign="top"
								  >
									<table
									  align="center"
									  border="0"
									  cellpadding="0"
									  cellspacing="0"
									  width="100%"
									  role="presentation"
									  data-block-id="14"
									>
									  <tbody>
										<tr class="mceRow">
										  <td
											style="
											  background-position: center;
											  background-repeat: no-repeat;
											  background-size: cover;
											"
											valign="top"
										  >
											<table border="0" cellpadding="0" cellspacing="0" width="100%" role="presentation">
											  <tbody>
												<tr>
												  <td
													style="padding-top: 0; padding-bottom: 0"
													class="mceColumn"
													data-block-id="-4"
													valign="top"
													colspan="12"
													width="100%"
												  >
													<table
													  border="0"
													  cellpadding="0"
													  cellspacing="0"
													  width="100%"
													  role="presentation"
													>
													  <tbody>
														<tr>
														  <td
															style="
															  padding-top: 12px;
															  padding-bottom: 12px;
															  padding-right: 48px;
															  padding-left: 48px;
															"
															class="mceBlockContainer"
															align="center"
															valign="top"
														  >
															<img
															  data-block-id="3"
															  width="150"
															  height="auto"
															  style="
																width: 150px;
																height: auto;
																max-width: 150px !important;
																display: block;
															  "
															  alt="Logo"
															  src="https://mcusercontent.com/beb0655d41009e40b7fb5a34d/images/86364ab7-2313-6890-abee-3d6095ca61ef.png"
															  class="mceLogo"
															/>
														  </td>
														</tr>
														<tr>
														  <td
															style="
															  padding-top: 12px;
															  padding-bottom: 12px;
															  padding-right: 24px;
															  padding-left: 24px;
															"
															class="mceBlockContainer"
															valign="top"
														  >
															<div
															  data-block-id="5"
															  class="mceText"
															  id="dataBlockId-5"
															  style="width: 100%"
															>
															  <h1 class="last-child" style="text-align: center">
																<span style="font-size: 30px">Log in to Predictmarket</span>
															  </h1>
															</div>
														  </td>
														</tr>
														<tr>
														  <td
															style="
															  padding-top: 12px;
															  padding-bottom: 12px;
															  padding-right: 100px;
															  padding-left: 100px;
															"
															class="mceBlockContainer"
															valign="top"
														  >
															<div
															  data-block-id="6"
															  class="mceText"
															  id="dataBlockId-6"
															  style="width: 100%"
															>
															  <p style="text-align: left">Hi ` + username + `,</p>
															  <p style="text-align: left">
																Copy and paste this temporary login code:
															  </p>
															</div>
														  </td>
														</tr>
														<tr>
														  <td
															style="
															  padding-top: 30px;
															  padding-bottom: 12px;
															  padding-right: 0;
															  padding-left: 0;
															"
															class="mceBlockContainer"
															align="center"
															valign="top"
														  >
															<table
															  align="center"
															  border="0"
															  cellpadding="0"
															  cellspacing="0"
															  role="presentation"
															  data-block-id="7"
															>
															  <tbody>
																<tr class="mceStandardButton">
																  <td
																	style="
																	  background-color: #000000;
																	  border-radius: 4px;
																	  text-align: center;
																	"
																	class="mceButton"
																	valign="top"
																  >
																	<p
																	  style="
																		background-color: #000000;
																		border-radius: 4px;
																		border: 2px solid #000000;
																		color: #ffffff;
																		display: block;
																		font-family: 'Helvetica Neue', Helvetica, Arial, Verdana,
																		  sans-serif;
																		font-size: 16px;
																		font-weight: normal;
																		font-style: normal;
																		padding: 5px 25px;
																		text-decoration: none;
																		min-width: 30px;
																		text-align: center;
																		direction: ltr;
																		letter-spacing: 0px;
																	  "
																	>
																	  ` + code + `
																	</p>
																  </td>
																</tr>
																<tr></tr>
															  </tbody>
															</table>
														  </td>
														</tr>
														<tr>
														  <td
															style="
															  padding-top: 30px;
															  padding-bottom: 12px;
															  padding-right: 100px;
															  padding-left: 100px;
															"
															class="mceBlockContainer"
															valign="top"
														  >
															<div
															  data-block-id="6"
															  class="mceText"
															  id="dataBlockId-6"
															  style="width: 100%"
															>
															  <p style="text-align: left">
																If you didn't try to login, you can safely ignore this email.
															  </p>
															</div>
														  </td>
														</tr>
														<tr>
														  <td
															style="
															  padding-top: 8px;
															  padding-bottom: 8px;
															  padding-right: 8px;
															  padding-left: 8px;
															"
															class="mceLayoutContainer"
															valign="top"
														  >
															<table
															  align="center"
															  border="0"
															  cellpadding="0"
															  cellspacing="0"
															  width="100%"
															  role="presentation"
															  data-block-id="13"
															  id="section_37951e71c1cf228f036a7cc5ce1086cc"
															  class="mceFooterSection"
															>
															  <tbody>
																<tr class="mceRow">
																  <td
																	style="
																	  background-position: center;
																	  background-repeat: no-repeat;
																	  background-size: cover;
																	"
																	valign="top"
																  >
																	<table
																	  border="0"
																	  cellpadding="0"
																	  cellspacing="12"
																	  width="100%"
																	  role="presentation"
																	>
																	  <tbody>
																		<tr>
																		  <td
																			style="
																			  padding-top: 0;
																			  padding-bottom: 0;
																			  margin-bottom: 12px;
																			"
																			class="mceColumn"
																			data-block-id="-3"
																			valign="top"
																			colspan="12"
																			width="100%"
																		  >
																			<table
																			  border="0"
																			  cellpadding="0"
																			  cellspacing="0"
																			  width="100%"
																			  role="presentation"
																			>
																			  <tbody>
																				<tr>
																				  <td
																					style="
																					  padding-top: 12px;
																					  padding-bottom: 12px;
																					  padding-right: 16px;
																					  padding-left: 16px;
																					"
																					class="mceBlockContainer"
																					align="center"
																					valign="top"
																				  >
																					<div
																					  data-block-id="11"
																					  class="mceText"
																					  id="dataBlockId-11"
																					  style="display: inline-block; width: 100%"
																					>
																					  <p class="last-child">
																						<em
																						  ><span style="font-size: 12px"
																							>@2024 PredictMarket Inc.</span
																						  ></em
																						>
																					  </p>
																					</div>
																				  </td>
																				</tr>
																			  </tbody>
																			</table>
																		  </td>
																		</tr>
																	  </tbody>
																	</table>
																  </td>
																</tr>
															  </tbody>
															</table>
														  </td>
														</tr>
													  </tbody>
													</table>
												  </td>
												</tr>
											  </tbody>
											</table>
										  </td>
										</tr>
									  </tbody>
									</table>
								  </td>
								</tr>
							  </tbody>
							</table>
						  </td>
						</tr>
					  </tbody>
					</table>
				  </td>
				</tr>
			  </tbody>
			</table>
		  </center>
		  ` +
			"</body></html>\r\n")
	return msg
}
