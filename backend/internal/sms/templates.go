package sms

type Template struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Body string `json:"body"`
}

var Templates = []Template{
	{
		Key:  "registration_approved",
		Name: "Başvuru onaylandı",
		Body: "KARE-REHBER sistemine hoş geldiniz. Telefon numaranız ve şu şifre ile giriş yapabilirsiniz: {{password}}",
	},
	{
		Key:  "user_created",
		Name: "Kullanıcı oluşturuldu",
		Body: "KARE-REHBER sistemine hoş geldiniz. Telefonunuz ve şu şifre ile giriş yapabilirsiniz: {{password}}",
	},
	{
		Key:  "password_reset",
		Name: "Şifre sıfırlandı",
		Body: "KARE-REHBER şifreniz sıfırlandı. Yeni şifreniz: {{password}}",
	},
	{
		Key:  "missing_evaluation_reminder",
		Name: "Eksik değerlendirme hatırlatması",
		Body: "Sayın {{name}}, {{week}} için öğrenci değerlendirmenizi henüz girmediniz. En kısa sürede tamamlamanızı rica ederiz.",
	},
	{
		Key:  "general",
		Name: "Genel mesaj",
		Body: "",
	},
}
