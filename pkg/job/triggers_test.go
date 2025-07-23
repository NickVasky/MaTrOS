package job

import (
	"reflect"
	"testing"
)

var yaml1 = `
triggers:
  firstTrigger:
    bot_id: Bot1
    process_id: Test_Process_1
    headers:
      - key: From
        value: me@gmail.com
    subject:
      - Test
      - Subject
    body:
      - Test

  secondTrigger:
    bot_id: BotFather
    process_id: Test_Process_2
    headers:
      - key: Date
        value: 2024-12-24
    body:
      - Робот
      - отчёт
`
var yaml1result = &TriggersConfigYaml{
	Triggers: map[string]TriggerYaml{
		"firstTrigger": TriggerYaml{
			BotID:     "Bot1",
			ProcessID: "Test_Process_1",
			Headers: []HeaderFieldYaml{
				{Key: "From", Value: "me@gmail.com"}},
			Subject: []string{"Test", "Subject"},
			Body:    []string{"Test"},
		},
		"secondTrigger": TriggerYaml{
			BotID:     "BotFather",
			ProcessID: "Test_Process_2",
			Headers: []HeaderFieldYaml{
				{Key: "Date", Value: "2024-12-24"}},
			Body: []string{"Робот", "отчёт"},
		},
	},
}

func Test_parseTriggersYaml(t *testing.T) {
	type args struct {
		yamlData []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *TriggersConfigYaml
		wantErr bool
	}{
		{
			name: "Test1",
			args: args{
				[]byte(yaml1),
			},
			want:    yaml1result,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTriggersYaml(tt.args.yamlData)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTriggersYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTriggersYaml() = %v, want %v", got, tt.want)
			}
		})
	}
}
