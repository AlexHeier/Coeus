package coeus

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var cons map[string]*Conversation

/*
Start starts the dashboard on the specified port.
The dashboard is a web interface for the LLM chatbot used for trubleshooting and testing the chatbot.

@param Port string - The port the dashboard should listen on.

@return error - Returns an error if the server could not start.
*/
func StartDashboard(port string) error {
	cons = make(map[string]*Conversation)
	http.HandleFunc("/api/chat", chatHandler)
	http.HandleFunc("/", webHandler)
	return http.ListenAndServe(":"+port, nil)
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		webGetHandler(w)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func webGetHandler(w http.ResponseWriter) {

	data, err := base64.StdEncoding.DecodeString(dashboardPageBase64)
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = w.Write(data)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		chatPostHandler(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func chatPostHandler(w http.ResponseWriter, r *http.Request) {
	req, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data map[string]interface{}

	err = json.Unmarshal(req, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, ok := data["userid"].(string)
	if !ok {
		http.Error(w, "bad userid type", http.StatusBadRequest)
		return
	}

	_, ok = data["prompt"].(string)
	if !ok {
		http.Error(w, "bad prompt type", http.StatusBadRequest)
		return
	}

	prompt := data["prompt"].(string)
	userid := data["userid"].(string)

	_, exist := cons[userid]
	if !exist {
		cons[userid] = BeginConversation()
	}

	res, err := cons[userid].Prompt(prompt)
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = w.Write([]byte(res.Response))
	if err != nil {
		fmt.Println(err.Error())
	}
}

// Base64 encoded HTML for the dashboard page
var dashboardPageBase64 string = "PCFET0NUWVBFIGh0bWw+DQo8aHRtbCBsYW5nPSJlbiI+DQo8aGVhZD4NCiAgICA8bWV0YSBjaGFyc2V0PSJVVEYtOCI+DQogICAgPG1ldGEgbmFtZT0idmlld3BvcnQiIGNvbnRlbnQ9IndpZHRoPWRldmljZS13aWR0aCwgaW5pdGlhbC1zY2FsZT0xLjAiPg0KICAgIDx0aXRsZT5Db2V1cyBEYXNoYm9hcmQ8L3RpdGxlPg0KICAgIDxzdHlsZT4NCiAgICAgICAgKiB7DQogICAgICAgICAgICBtYXJnaW46IDA7DQogICAgICAgICAgICBwYWRkaW5nOiAwOw0KICAgICAgICAgICAgYm94LXNpemluZzogYm9yZGVyLWJveDsNCiAgICAgICAgICAgIGZvbnQtZmFtaWx5OiBBcmlhbCwgc2Fucy1zZXJpZjsNCiAgICAgICAgfQ0KDQogICAgICAgIGJvZHkgew0KICAgICAgICAgICAgZGlzcGxheTogZmxleDsNCiAgICAgICAgICAgIGhlaWdodDogMTAwdmg7DQogICAgICAgIH0NCg0KICAgICAgICAvKiBTaWRlYmFyIFN0eWxpbmcgKi8NCiAgICAgICAgLnNpZGViYXIgew0KICAgICAgICAgICAgd2lkdGg6IDI1JTsNCiAgICAgICAgICAgIGJhY2tncm91bmQ6ICNmMGYwZjA7DQogICAgICAgICAgICBib3JkZXItcmlnaHQ6IDFweCBzb2xpZCAjY2NjOw0KICAgICAgICAgICAgb3ZlcmZsb3cteTogYXV0bzsNCiAgICAgICAgICAgIHBhZGRpbmc6IDEwcHg7DQogICAgICAgIH0NCg0KICAgICAgICAuY2hhdC1saXN0IHsNCiAgICAgICAgICAgIGxpc3Qtc3R5bGU6IG5vbmU7DQogICAgICAgIH0NCg0KICAgICAgICAuY2hhdC1saXN0IGxpIHsNCiAgICAgICAgICAgIHBhZGRpbmc6IDE1cHg7DQogICAgICAgICAgICBjdXJzb3I6IHBvaW50ZXI7DQogICAgICAgICAgICBib3JkZXItYm90dG9tOiAxcHggc29saWQgI2RkZDsNCiAgICAgICAgICAgIHRyYW5zaXRpb246IGJhY2tncm91bmQgMC4zczsNCiAgICAgICAgfQ0KDQogICAgICAgIC5jaGF0LWxpc3QgbGk6aG92ZXIsIC5jaGF0LWxpc3QgLmFjdGl2ZSB7DQogICAgICAgICAgICBiYWNrZ3JvdW5kOiAjZDlkOWQ5Ow0KICAgICAgICB9DQoNCiAgICAgICAgLyogQ2hhdCBXaW5kb3cgU3R5bGluZyAqLw0KICAgICAgICAuY2hhdC1jb250YWluZXIgew0KICAgICAgICAgICAgd2lkdGg6IDc1JTsNCiAgICAgICAgICAgIGRpc3BsYXk6IGZsZXg7DQogICAgICAgICAgICBmbGV4LWRpcmVjdGlvbjogY29sdW1uOw0KICAgICAgICAgICAgYmFja2dyb3VuZDogI2ZmZjsNCiAgICAgICAgfQ0KDQogICAgICAgIC5jaGF0LWhlYWRlciB7DQogICAgICAgICAgICBwYWRkaW5nOiAxNXB4Ow0KICAgICAgICAgICAgYm9yZGVyLWJvdHRvbTogMXB4IHNvbGlkICNjY2M7DQogICAgICAgICAgICBiYWNrZ3JvdW5kOiAjZjBmMGYwOw0KICAgICAgICAgICAgZm9udC13ZWlnaHQ6IGJvbGQ7DQogICAgICAgIH0NCg0KICAgICAgICAuY2hhdC1ib3ggew0KICAgICAgICAgICAgZmxleDogMTsNCiAgICAgICAgICAgIHBhZGRpbmc6IDE1cHg7DQogICAgICAgICAgICBvdmVyZmxvdy15OiBhdXRvOw0KICAgICAgICAgICAgZGlzcGxheTogZmxleDsNCiAgICAgICAgICAgIGZsZXgtZGlyZWN0aW9uOiBjb2x1bW47DQogICAgICAgIH0NCg0KICAgICAgICAubWVzc2FnZSB7DQogICAgICAgICAgICBtYXgtd2lkdGg6IDcwJTsNCiAgICAgICAgICAgIHBhZGRpbmc6IDEwcHg7DQogICAgICAgICAgICBtYXJnaW4tYm90dG9tOiAxMHB4Ow0KICAgICAgICAgICAgYm9yZGVyLXJhZGl1czogMTBweDsNCiAgICAgICAgfQ0KDQogICAgICAgIC5zZW50IHsNCiAgICAgICAgICAgIGFsaWduLXNlbGY6IGZsZXgtZW5kOw0KICAgICAgICAgICAgYmFja2dyb3VuZDogIzAwNzhmZjsNCiAgICAgICAgICAgIGNvbG9yOiB3aGl0ZTsNCiAgICAgICAgfQ0KDQogICAgICAgIC5yZWNlaXZlZCB7DQogICAgICAgICAgICBhbGlnbi1zZWxmOiBmbGV4LXN0YXJ0Ow0KICAgICAgICAgICAgYmFja2dyb3VuZDogI2UwZTBlMDsNCiAgICAgICAgfQ0KDQogICAgICAgIC8qIE1lc3NhZ2UgSW5wdXQgQm94ICovDQogICAgICAgIC5jaGF0LWlucHV0IHsNCiAgICAgICAgICAgIGRpc3BsYXk6IGZsZXg7DQogICAgICAgICAgICBib3JkZXItdG9wOiAxcHggc29saWQgI2NjYzsNCiAgICAgICAgICAgIHBhZGRpbmc6IDEwcHg7DQogICAgICAgICAgICBiYWNrZ3JvdW5kOiAjZjlmOWY5Ow0KICAgICAgICB9DQoNCiAgICAgICAgLmNoYXQtaW5wdXQgaW5wdXQgew0KICAgICAgICAgICAgZmxleDogMTsNCiAgICAgICAgICAgIHBhZGRpbmc6IDEwcHg7DQogICAgICAgICAgICBib3JkZXI6IG5vbmU7DQogICAgICAgICAgICBib3JkZXItcmFkaXVzOiA1cHg7DQogICAgICAgICAgICBvdXRsaW5lOiBub25lOw0KICAgICAgICB9DQoNCiAgICAgICAgYnV0dG9uIHsNCiAgICAgICAgICAgIHBhZGRpbmc6IDEwcHggMTBweDsNCiAgICAgICAgICAgIGJhY2tncm91bmQ6ICMwMDc4ZmY7DQogICAgICAgICAgICBjb2xvcjogd2hpdGU7DQogICAgICAgICAgICBib3JkZXI6IG5vbmU7DQogICAgICAgICAgICBib3JkZXItcmFkaXVzOiA1cHg7DQogICAgICAgICAgICBjdXJzb3I6IHBvaW50ZXI7DQogICAgICAgICAgICBtYXJnaW4tbGVmdDogNXB4Ow0KICAgICAgICAgICAgbWFyZ2luLXJpZ2h0OiA1cHg7DQogICAgICAgIH0NCg0KICAgICAgICAjc2VydmVyLWJveCB7DQogICAgICAgICAgICBiYWNrZ3JvdW5kOiAjMDA3OGZmOw0KICAgICAgICAgICAgY29sb3I6IHdoaXRlOw0KICAgICAgICAgICAgYm9yZGVyOiBub25lOw0KICAgICAgICAgICAgYm9yZGVyLXJhZGl1czogNXB4Ow0KICAgICAgICAgICAgbWFyZ2luOiA1cHg7DQogICAgICAgICAgICBwYWRkaW5nOiA1cHg7DQogICAgICAgICAgICBhbGlnbi1jb250ZW50OiBjZW50ZXI7DQogICAgICAgICAgICB3aWR0aDogOTYlOw0KICAgICAgICB9DQoNCiAgICAgICAgI3NlcnZlci1ib3ggcCB7DQogICAgICAgICAgICB0ZXh0LWFsaWduOiBjZW50ZXI7DQogICAgICAgIH0NCg0KICAgICAgICAjbmV3Q2hhdEJ1dHRvbiB7DQogICAgICAgICAgICB3aWR0aDogOTYlOw0KICAgICAgICB9DQoNCiAgICAgICAgI2lucHV0X3NlcnZlcl9pcCB7DQogICAgICAgICAgICB3aWR0aDogOTMlOw0KICAgICAgICAgICAgcGFkZGluZzogNXB4Ow0KICAgICAgICAgICAgbWFyZ2luLXRvcDogNXB4Ow0KICAgICAgICAgICAgbWFyZ2luLWxlZnQ6IDEwcHg7DQogICAgICAgICAgICBtYXJnaW4tYm90dG9tOiA1cHg7DQogICAgICAgICAgICBib3JkZXItcmFkaXVzOiA1cHg7DQogICAgICAgIH0NCg0KICAgICAgICAuc2VsZWN0b3Igew0KICAgICAgICAgICAgcGFkZGluZzogMTBweCAxNXB4Ow0KICAgICAgICAgICAgYmFja2dyb3VuZDogIzAwNzhmZjsNCiAgICAgICAgICAgIGNvbG9yOiB3aGl0ZTsNCiAgICAgICAgICAgIGJvcmRlcjogbm9uZTsNCiAgICAgICAgICAgIGJvcmRlci1yYWRpdXM6IDVweDsNCiAgICAgICAgICAgIGN1cnNvcjogcG9pbnRlcjsNCiAgICAgICAgICAgIG1hcmdpbi1sZWZ0OiAxMHB4Ow0KICAgICAgICAgICAgdGV4dC1vdmVyZmxvdzogZWxsaXBzaXM7DQogICAgICAgICAgICBtaW4td2lkdGg6IDEyMHB4Ow0KICAgICAgICAgICAgbWF4LXdpZHRoOiAxMCU7DQogICAgICAgIH0NCg0KICAgIDwvc3R5bGU+DQo8L2hlYWQ+DQo8Ym9keT4NCg0KICAgIDwhLS0gU2lkZWJhciBmb3IgQ2hhdCBTZWxlY3Rpb24gLS0+DQogICAgPGRpdiBjbGFzcz0ic2lkZWJhciI+DQogICAgICAgIDxkaXYgaWQ9InNlcnZlci1ib3giPg0KICAgICAgICAgICAgPHA+Q29ldXMgU2VydmVyPC9wPg0KICAgICAgICAgICAgPGlucHV0IHR5cGU9InRleHQiIHZhbHVlPSIiICBpZD0iaW5wdXRfc2VydmVyX2lwIj4NCiAgICAgICAgPC9kaXY+DQogICAgICAgIDxidXR0b24gaWQ9Im5ld0NoYXRCdXR0b24iIG9uY2xpY2s9ImNyZWF0ZUNvbnZlcnNhdGlvbigpIj4NCiAgICAgICAgICAgIE5ld0NoYXQNCiAgICAgICAgPC9idXR0b24+DQoNCg0KICAgICAgICA8dWwgY2xhc3M9ImNoYXQtbGlzdCI+DQogICAgICAgIDwvdWw+DQogICAgPC9kaXY+DQoNCiAgICA8IS0tIENoYXQgV2luZG93IC0tPg0KICAgIDxkaXYgY2xhc3M9ImNoYXQtY29udGFpbmVyIj4NCiAgICAgICAgPGRpdiBjbGFzcz0iY2hhdC1oZWFkZXIiPkNyZWF0ZSBhIG5ldyBjaGF0PC9kaXY+DQogICAgICAgIDxkaXYgY2xhc3M9ImNoYXQtYm94Ij48L2Rpdj4NCiAgICAgICAgPGRpdiBjbGFzcz0iY2hhdC1pbnB1dCI+DQogICAgICAgICAgICA8aW5wdXQgdHlwZT0idGV4dCIgaWQ9Im1lc3NhZ2VJbnB1dCIgcGxhY2Vob2xkZXI9IlR5cGUgYSBtZXNzYWdlLi4uIj4NCiAgICAgICAgICAgIDxidXR0b24gb25jbGljaz0iTWVzc2FnZUxMTSgpIj5TZW5kPC9idXR0b24+DQogICAgICAgICAgICA8YnV0dG9uIGlkPSJzcGVlY2hCdXR0b24iPlNwZWFrPC9idXR0b24+DQogICAgICAgICAgICA8c2VsZWN0IGNsYXNzPSJzZWxlY3RvciIgbmFtZT0ibGFuZ3VhZ2VTZWxlY3RvciIgaWQ9Imxhbmd1YWdlU2VsZWN0b3IiPg0KICAgICAgICAgICAgICAgIDxvcHRpb24gdmFsdWU9ImVuLVVTIj5TcGVlY2ggVG8gVGV4dCBMYW5ndWFnZTwvb3B0aW9uPg0KICAgICAgICAgICAgICAgIDxvcHRpb24gdmFsdWU9ImVuLVVTIj5FbmdsaXNoIFVTPC9vcHRpb24+DQogICAgICAgICAgICAgICAgPG9wdGlvbiB2YWx1ZT0iZW4tR0IiPkVuZ2xpc2ggVUs8L29wdGlvbj4NCiAgICAgICAgICAgICAgICA8b3B0aW9uIHZhbHVlPSJuby1OQiI+Tm9yd2VnaWFuPC9vcHRpb24+DQogICAgICAgICAgICA8L3NlbGVjdD4NCiAgICAgICAgICAgIDxzZWxlY3QgY2xhc3M9InNlbGVjdG9yIiBuYW1lPSJ2b2ljZVNlbGVjdG9yIiBpZD0idm9pY2VTZWxlY3RvciI+DQogICAgICAgICAgICAgICAgPG9wdGlvbiB2YWx1ZT0iIj5ObyBUZXh0IFRvIFNwZWVjaDwvb3B0aW9uPg0KICAgICAgICAgICAgPC9zZWxlY3Q+DQogICAgICAgIDwvZGl2Pg0KICAgIDwvZGl2Pg0KDQogICAgPHNjcmlwdD4NCg0KICAgICAgICBsZXQgY29udmVyc2F0aW9uczsgLy8gU3RvcmVzIHRoZSBjb252ZXJzYXRpb25zDQogICAgICAgIGxldCBjdXJyZW50Q2hhdDsgLy8gUG9pbnRzIHRvIHRoZSBhY3RpdmUgY29udmVyc2F0aW9uDQoNCiAgICAgICAgY29uc3Qgc3RhcnRTcGVhY2ggPSBkb2N1bWVudC5nZXRFbGVtZW50QnlJZCgic3BlZWNoQnV0dG9uIik7DQogICAgICAgIGNvbnN0IG1lc3NhZ2VJbnB1dCA9IGRvY3VtZW50LmdldEVsZW1lbnRCeUlkKCJtZXNzYWdlSW5wdXQiKTsNCiAgICAgICAgY29uc3QgbGFuZ3VhZ2VTZWxlY3QgPSBkb2N1bWVudC5nZXRFbGVtZW50QnlJZCgibGFuZ3VhZ2VTZWxlY3RvciIpOw0KICAgICAgIA0KICAgICAgICBjb25zdCBTcGVlY2hSZWNvZ25pdGlvbiA9IHdpbmRvdy5TcGVlY2hSZWNvZ25pdGlvbiB8fCB3aW5kb3cud2Via2l0U3BlZWNoUmVjb2duaXRpb247IC8vIFVzZWQgaW4gY29udmVydGluZyBzcGVlY2ggdG8gdGV4dA0KICAgICAgICBjb25zdCBzeW50aCA9IHdpbmRvdy5zcGVlY2hTeW50aGVzaXM7IC8vIFVzZWQgaW4gY29udmVydGluZyB0ZXh0IHRvIHNwZWVjaA0KICAgICAgICBsZXQgdm9pY2VzID0gc3ludGguZ2V0Vm9pY2VzKCk7IC8vIFN0b3JlcyB0aGUgdm9pY2VzIHVzZWQgaW4gdGV4dCB0byBzcGVlY2ggc3ludGgNCiAgICAgICAgDQogICAgICAgIGRvY3VtZW50LmdldEVsZW1lbnRCeUlkKCJpbnB1dF9zZXJ2ZXJfaXAiKS52YWx1ZSA9IHdpbmRvdy5sb2NhdGlvbi5ob3N0Ow0KDQogICAgICAgIGlmIChTcGVlY2hSZWNvZ25pdGlvbikgew0KDQogICAgICAgICAgICBjb25zdCByZWNvZ25pdGlvbiA9IG5ldyBTcGVlY2hSZWNvZ25pdGlvbigpOw0KDQogICAgICAgICAgICByZWNvZ25pdGlvbi5jb250aW51b3VzID0gZmFsc2U7IC8vIFN0b3Agd2hlbiBzcGVlY2ggZW5kcw0KICAgICAgICAgICAgcmVjb2duaXRpb24ubGFuZyA9ICdlbi1VUyc7IC8vIExhbmd1YWdlDQogICAgICAgICAgICByZWNvZ25pdGlvbi5pbnRlcmltUmVzdWx0cyA9IHRydWU7IC8vIFNob3cgaW50ZXJpbSByZXN1bHRzDQoNCiAgICAgICAgICAgIC8vIFN0YXJ0IHJlY29nbml0aW9uDQogICAgICAgICAgICBzdGFydFNwZWFjaC5hZGRFdmVudExpc3RlbmVyKCdjbGljaycsICgpID0+IHsNCiAgICAgICAgICAgIHJlY29nbml0aW9uLnN0YXJ0KCk7DQogICAgICAgICAgICB9KTsNCg0KICAgICAgICAgICAgbGFuZ3VhZ2VTZWxlY3QuYWRkRXZlbnRMaXN0ZW5lcignY2hhbmdlJywgKCkgPT4gew0KICAgICAgICAgICAgICAgIHJlY29nbml0aW9uLmxhbmcgPSBsYW5ndWFnZVNlbGVjdC52YWx1ZTsNCiAgICAgICAgICAgIH0pOw0KDQogICAgICAgICAgICAvLyBDYXB0dXJlIHRoZSByZXN1bHQNCiAgICAgICAgICAgIHJlY29nbml0aW9uLm9ucmVzdWx0ID0gKGV2ZW50KSA9PiB7DQogICAgICAgICAgICAgICAgbGV0IHRyYW5zY3JpcHQgPSAnJzsNCiAgICAgICAgICAgICAgICBmb3IgKGxldCBpID0gMDsgaSA8IGV2ZW50LnJlc3VsdHMubGVuZ3RoOyBpKyspIHsNCiAgICAgICAgICAgICAgICAgICAgdHJhbnNjcmlwdCArPSBldmVudC5yZXN1bHRzW2ldWzBdLnRyYW5zY3JpcHQ7DQogICAgICAgICAgICAgICAgfQ0KICAgICAgICAgICAgICAgIG1lc3NhZ2VJbnB1dC52YWx1ZSA9IHRyYW5zY3JpcHQ7DQogICAgICAgICAgICAgICAgLy90ZXh0ID0gdHJhbnNjcmlwdDsNCiAgICAgICAgICAgIH07DQoNCiAgICAgICAgICAgIC8vIEFjdGlvbnMgdG8gZG8gd2hlbiB0aGUgY2FwdHVyZSBlbmRzDQogICAgICAgICAgICByZWNvZ25pdGlvbi5vbmVuZCA9IChldmVudCkgPT4gew0KICAgICAgICAgICAgICAgIE1lc3NhZ2VMTE0oKTsNCiAgICAgICAgICAgIH0NCg0KICAgICAgICAgICAgLy8gRXJyb3IgaGFuZGxpbmcNCiAgICAgICAgICAgIHJlY29nbml0aW9uLm9uZXJyb3IgPSAoZXZlbnQpID0+IHsNCiAgICAgICAgICAgICAgICBtZXNzYWdlSW5wdXQudmFsdWUgPSAnRXJyb3Igb2NjdXJyZWQ6ICcgKyBldmVudC5lcnJvcjsNCiAgICAgICAgICAgIH07DQogICAgICAgIH0gZWxzZSB7DQogICAgICAgICAgICBzdGFydFNwZWFjaC5zdHlsZS5iYWNrZ3JvdW5kQ29sb3IgPSAiZ3JleSI7DQogICAgICAgIH07DQoNCg0KICAgICAgICBmdW5jdGlvbiBwb3B1bGF0ZVZvaWNlTGlzdCgpIHsNCiAgICAgICAgICAgIGNvbnN0IHZvaWNlU2VsZWN0ID0gZG9jdW1lbnQuZ2V0RWxlbWVudEJ5SWQoInZvaWNlU2VsZWN0b3IiKTsNCg0KICAgICAgICAgICAgdm9pY2VzLmZvckVhY2godm9pY2UgPT4gew0KICAgICAgICAgICAgICAgIGNvbnN0IG9wdGlvbiA9IGRvY3VtZW50LmNyZWF0ZUVsZW1lbnQoIm9wdGlvbiIpOw0KICAgICAgICAgICAgICAgIG9wdGlvbi50ZXh0Q29udGVudCA9IGAke3ZvaWNlLm5hbWV9ICgke3ZvaWNlLmxhbmd9KWA7DQoNCiAgICAgICAgICAgICAgICBpZiAodm9pY2UuZGVmYXVsdCkgew0KICAgICAgICAgICAgICAgICAgICBvcHRpb24udGV4dENvbnRlbnQgKz0gIiDigJQgREVGQVVMVCI7DQogICAgICAgICAgICAgICAgfQ0KDQogICAgICAgICAgICAgICAgb3B0aW9uLnNldEF0dHJpYnV0ZSgiZGF0YS1sYW5nIiwgdm9pY2UubGFuZyk7DQogICAgICAgICAgICAgICAgb3B0aW9uLnNldEF0dHJpYnV0ZSgiZGF0YS1uYW1lIiwgdm9pY2UubmFtZSk7DQogICAgICAgICAgICAgICAgdm9pY2VTZWxlY3QuYXBwZW5kQ2hpbGQob3B0aW9uKTsNCiAgICAgICAgICAgIH0pOw0KICAgICAgICB9Ow0KDQogICAgICAgIHdpbmRvdy5zcGVlY2hTeW50aGVzaXMub252b2ljZXNjaGFuZ2VkID0gZnVuY3Rpb24oKSB7DQogICAgICAgICAgICB2b2ljZXMgPSBzeW50aC5nZXRWb2ljZXMoKTsNCiAgICAgICAgICAgIHBvcHVsYXRlVm9pY2VMaXN0KCk7DQogICAgICAgIH07ICANCg0KICAgICAgICAvLyBGdW5jdGlvbiB0byBzd2l0Y2ggY2hhdHMNCiAgICAgICAgZnVuY3Rpb24gc3dpdGNoQ2hhdChpbmRleCkgew0KICAgICAgICAgICAgY3VycmVudENoYXQgPSBpbmRleDsNCg0KICAgICAgICAgICAgLy8gVXBkYXRlIGFjdGl2ZSBjbGFzcw0KICAgICAgICAgICAgZG9jdW1lbnQucXVlcnlTZWxlY3RvckFsbCgiLmNoYXQtbGlzdCBsaSIpLmZvckVhY2goKGxpLCBpKSA9PiB7DQogICAgICAgICAgICAgICAgaWYgKGxpLmlubmVySFRNTCA9PSBgQ2hhdC0ke2N1cnJlbnRDaGF0fWApIHsNCiAgICAgICAgICAgICAgICAgICAgbGkuY2xhc3NMaXN0LmFkZCgiYWN0aXZlIik7DQogICAgICAgICAgICAgICAgfSBlbHNlIHsNCiAgICAgICAgICAgICAgICAgICAgbGkuY2xhc3NMaXN0LnJlbW92ZSgiYWN0aXZlIik7DQogICAgICAgICAgICAgICAgfQ0KICAgICAgICAgICAgfSk7DQoNCiAgICAgICAgICAgIGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3IoIi5jaGF0LWhlYWRlciIpLnRleHRDb250ZW50ID0gY3VycmVudENoYXQ7DQogICAgICAgICAgICBsb2FkTWVzc2FnZXMoKTsNCiAgICAgICAgfTsNCg0KICAgICAgICAvLyBGdW5jdGlvbiB0byBsb2FkIG1lc3NhZ2VzDQogICAgICAgIGZ1bmN0aW9uIGxvYWRNZXNzYWdlcygpIHsNCiAgICAgICAgICAgIGNvbnN0IGNoYXRCb3ggPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yKCIuY2hhdC1ib3giKTsNCiAgICAgICAgICAgIGNoYXRCb3guaW5uZXJIVE1MID0gIiI7DQoNCiAgICAgICAgICAgIGNvbnZlcnNhdGlvbnMuZ2V0KGN1cnJlbnRDaGF0KS5mb3JFYWNoKG1zZyA9PiB7DQogICAgICAgICAgICAgICAgY29uc3QgbWVzc2FnZURpdiA9IGRvY3VtZW50LmNyZWF0ZUVsZW1lbnQoImRpdiIpOw0KICAgICAgICAgICAgICAgIG1lc3NhZ2VEaXYuY2xhc3NMaXN0LmFkZCgibWVzc2FnZSIsIG1zZy50eXBlKTsNCiAgICAgICAgICAgICAgICBtZXNzYWdlRGl2LnRleHRDb250ZW50ID0gbXNnLnRleHQ7DQogICAgICAgICAgICAgICAgY2hhdEJveC5hcHBlbmRDaGlsZChtZXNzYWdlRGl2KTsNCiAgICAgICAgICAgIH0pOw0KDQogICAgICAgICAgICBjaGF0Qm94LnNjcm9sbFRvcCA9IGNoYXRCb3guc2Nyb2xsSGVpZ2h0Ow0KICAgICAgICB9Ow0KDQogICAgICAgIGZ1bmN0aW9uIGNyZWF0ZUNvbnZlcnNhdGlvbigpIHsNCiAgICAgICAgICAgIHNldEN1cnJlbnRDaGF0KCkNCiAgICAgICAgICAgIHZhciBjaGF0bGlzdCA9IGRvY3VtZW50LmdldEVsZW1lbnRzQnlDbGFzc05hbWUoImNoYXQtbGlzdCIpWzBdOw0KICAgICAgICAgICAgdmFyIGNvbnZlcnNhdGlvbiA9IGRvY3VtZW50LmNyZWF0ZUVsZW1lbnQoImxpIik7DQogICAgICAgICAgICBjb252ZXJzYXRpb24uc2V0QXR0cmlidXRlKCdvbmNsaWNrJywgYHN3aXRjaENoYXQoJyR7Y3VycmVudENoYXR9JylgKTsNCiAgICAgICAgICAgIGNvbnZlcnNhdGlvbi5pbm5lckhUTUwgPSBgQ2hhdC0ke2N1cnJlbnRDaGF0fWA7DQogICAgICAgICAgICBjaGF0bGlzdC5hcHBlbmRDaGlsZChjb252ZXJzYXRpb24pOw0KICAgICAgICAgICAgY29udmVyc2F0aW9ucy5zZXQoY3VycmVudENoYXQsIFt7IHRleHQ6ICJOZXcgY29udmVyc2F0aW9uLiBTZW5kIGEgbWVzc2FnZSB0byBiZWdpbiB0aGUgY29udmVyc2F0aW9uIiwgdHlwZTogInJlY2VpdmVkIiB9XSk7DQogICAgICAgICAgICBkb2N1bWVudC5xdWVyeVNlbGVjdG9yKCIuY2hhdC1oZWFkZXIiKS50ZXh0Q29udGVudCA9IGN1cnJlbnRDaGF0Ow0KICAgICAgICB9Ow0KDQogICAgICAgIGZ1bmN0aW9uIE1lc3NhZ2VMTE0oKSB7DQoNCiAgICAgICAgICAgIHVybCA9ICJodHRwOi8vIiArIGAke2RvY3VtZW50LmdldEVsZW1lbnRCeUlkKCJpbnB1dF9zZXJ2ZXJfaXAiKS52YWx1ZX1gKyAiL2FwaS9jaGF0Ig0KDQogICAgICAgICAgICBjb25zdCBpbnB1dCA9IGRvY3VtZW50LmdldEVsZW1lbnRCeUlkKCJtZXNzYWdlSW5wdXQiKTsNCiAgICAgICAgICAgIGNvbnN0IHRleHQgPSBpbnB1dC52YWx1ZS50cmltKCk7DQogICAgICAgICAgICBpZiAodGV4dCA9PT0gIiIpIHJldHVybjsNCg0KICAgICAgICAgICAgLy8gQWRkIG1lc3NhZ2UgdG8gY29udmVyc2F0aW9uDQogICAgICAgICAgICBjb252ZXJzYXRpb25zLmdldChjdXJyZW50Q2hhdCkucHVzaCh7IHRleHQsIHR5cGU6ICJzZW50IiB9KTsNCiAgICAgICAgICAgIGxvYWRNZXNzYWdlcygpOw0KDQogICAgICAgICAgICAvLyBDbGVhciBpbnB1dA0KICAgICAgICAgICAgaW5wdXQudmFsdWUgPSAiIjsNCg0KICAgICAgICAgICAgYm9keSA9IEpTT04uc3RyaW5naWZ5KHsNCiAgICAgICAgICAgICAgICAgICAgdXNlcmlkOiBjdXJyZW50Q2hhdCwNCiAgICAgICAgICAgICAgICAgICAgcHJvbXB0OiB0ZXh0DQogICAgICAgICAgICAgICAgfSk7DQoNCiAgICAgICAgICAgIGNvbnN0IHJlc3BvbnNlID0gZmV0Y2godXJsLCB7DQogICAgICAgICAgICAgICAgbWV0aG9kOiAnUE9TVCcsDQogICAgICAgICAgICAgICAgYm9keTogYm9keSwNCiAgICAgICAgICAgICAgICAgaGVhZGVyczogew0KICAgICAgICAgICAgICAgICAgICAnQWNjZXB0JzogJ2FwcGxpY2F0aW9uL2pzb24nLA0KICAgICAgICAgICAgICAgICAgICAnQ29udGVudC1UeXBlJzogJ2FwcGxpY2F0aW9uL2pzb24nLA0KICAgICAgICAgICAgICAgIH0sDQogICAgICAgICAgICB9KQ0KICAgICAgICAgICAgLnRoZW4ocmVzcG9uc2UgPT4gcmVzcG9uc2UudGV4dCgpKQ0KICAgICAgICAgICAgLnRoZW4oZGF0YSA9PiB7IA0KICAgICAgICAgICAgICAgIHNldFRpbWVvdXQoKCkgPT4gew0KICAgICAgICAgICAgICAgICAgICBjb252ZXJzYXRpb25zLmdldChjdXJyZW50Q2hhdCkucHVzaCh7IHRleHQ6IGRhdGEsIHR5cGU6ICJyZWNlaXZlZCIgfSk7DQogICAgICAgICAgICAgICAgICAgIGxvYWRNZXNzYWdlcygpOw0KDQogICAgICAgICAgICAgICAgICAgIC8vIElmIGEgdm9pY2UgaXMgc2VsZWN0ZWQgdGhlbiBzcGVhayBvdXQgdGhlIExMTSByZXNwb25zZQ0KICAgICAgICAgICAgICAgICAgICBpZiAoZG9jdW1lbnQuZ2V0RWxlbWVudEJ5SWQoInZvaWNlU2VsZWN0b3IiKS52YWx1ZSAhPT0gIiIpIHsNCiAgICAgICAgICAgICAgICAgICAgICAgIGNvbnN0IHJlc3BvbnNlID0gbmV3IFNwZWVjaFN5bnRoZXNpc1V0dGVyYW5jZShkYXRhKTsNCiAgICAgICAgICAgICAgICAgICAgICAgIGNvbnN0IHNlbGVjdGVkT3B0aW9uID0gZG9jdW1lbnQuZ2V0RWxlbWVudEJ5SWQoInZvaWNlU2VsZWN0b3IiKS5zZWxlY3RlZE9wdGlvbnNbMF0uZ2V0QXR0cmlidXRlKCJkYXRhLW5hbWUiKTsNCg0KICAgICAgICAgICAgICAgICAgICAgICAgdm9pY2VzLmZvckVhY2godm9pY2UgPT4geyANCiAgICAgICAgICAgICAgICAgICAgICAgICAgICBpZiAodm9pY2UubmFtZSA9PT0gc2VsZWN0ZWRPcHRpb24pIHsNCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgcmVzcG9uc2Uudm9pY2UgPSB2b2ljZTsNCiAgICAgICAgICAgICAgICAgICAgICAgICAgICB9DQogICAgICAgICAgICAgICAgICAgICAgICB9KTsNCg0KICAgICAgICAgICAgICAgICAgICAgICAgcmVzcG9uc2UucGl0Y2ggPSAxOw0KICAgICAgICAgICAgICAgICAgICAgICAgcmVzcG9uc2UucmF0ZSA9IDE7DQogICAgICAgICAgICAgICAgICAgICAgICBzeW50aC5zcGVhayhyZXNwb25zZSk7DQogICAgICAgICAgICAgICAgICAgIH0NCiAgICAgICAgICAgICAgICAgfSwgMTAwMCk7DQogICAgICAgICAgICB9KTsNCiAgICAgICAgIH07DQoNCiAgICAgICAgZnVuY3Rpb24gc2V0Q3VycmVudENoYXQoKSB7DQogICAgICAgICAgICBjdXJyZW50Q2hhdCA9IChgJHtNYXRoLmZsb29yKE1hdGgucmFuZG9tKCkgKiAxMDAwMCl9YCkudG9TdHJpbmcoKTsNCiAgICAgICAgfTsNCg0KICAgICAgICBmdW5jdGlvbiBzZXR1cE1hcCgpIHsNCiAgICAgICAgICAgIGNvbnZlcnNhdGlvbnMgPSBuZXcgTWFwKCk7DQogICAgICAgIH07DQogICAgICAgIA0KDQoNCiAgICAgICAgc2V0dXBNYXAoKTsNCiAgICAgICAgc2V0Q3VycmVudENoYXQoKTsNCiAgICAgICAgY3JlYXRlQ29udmVyc2F0aW9uKGN1cnJlbnRDaGF0KTsNCiAgICAgICAgbG9hZE1lc3NhZ2VzKCk7DQoNCiAgICAgICAgZG9jdW1lbnQucXVlcnlTZWxlY3RvcigiI21lc3NhZ2VJbnB1dCIpLmFkZEV2ZW50TGlzdGVuZXIoImtleXVwIiwgZXZlbnQgPT4gew0KICAgICAgICAgICAgaWYoZXZlbnQua2V5ICE9PSAiRW50ZXIiKSByZXR1cm47IC8vIFVzZSBgLmtleWAgaW5zdGVhZC4NCiAgICAgICAgICAgIE1lc3NhZ2VMTE0oKSAvLyBUaGluZ3MgeW91IHdhbnQgdG8gZG8uDQogICAgICAgICAgICBldmVudC5wcmV2ZW50RGVmYXVsdCgpOyAvLyBObyBuZWVkIHRvIGByZXR1cm4gZmFsc2U7YC4NCiAgICAgICAgfSk7DQogICAgPC9zY3JpcHQ+DQoNCjwvYm9keT4NCjwvaHRtbD4NCg=="
