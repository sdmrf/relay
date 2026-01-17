package runtime

import "testing"

func TestParseJavaVersion(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		want    int
		wantErr bool
	}{
		{
			name:   "Java 17",
			output: `openjdk version "17.0.8" 2023-07-18`,
			want:   17,
		},
		{
			name:   "Java 21",
			output: `openjdk version "21.0.1" 2023-10-17`,
			want:   21,
		},
		{
			name:   "Java 11",
			output: `openjdk version "11.0.20" 2023-07-18`,
			want:   11,
		},
		{
			name:   "Java 8 (1.8 format)",
			output: `java version "1.8.0_392"`,
			want:   8,
		},
		{
			name:   "Java 8 OpenJDK",
			output: `openjdk version "1.8.0_362"`,
			want:   8,
		},
		{
			name:   "Oracle Java 17",
			output: `java version "17.0.8" 2023-07-18 LTS`,
			want:   17,
		},
		{
			name:   "Temurin 21",
			output: `openjdk version "21" 2023-09-19`,
			want:   21,
		},
		{
			name:    "No quotes",
			output:  `openjdk version 17.0.8`,
			wantErr: true,
		},
		{
			name:    "Empty output",
			output:  ``,
			wantErr: true,
		},
		{
			name:    "Garbage",
			output:  `not java version output`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseJavaVersion(tt.output)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJavaVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseJavaVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseJavaVersionMultiline(t *testing.T) {
	output := `openjdk version "17.0.8" 2023-07-18
OpenJDK Runtime Environment Temurin-17.0.8+7 (build 17.0.8+7)
OpenJDK 64-Bit Server VM Temurin-17.0.8+7 (build 17.0.8+7, mixed mode)`

	got, err := ParseJavaVersion(output)
	if err != nil {
		t.Fatalf("ParseJavaVersion() error = %v", err)
	}

	if got != 17 {
		t.Errorf("ParseJavaVersion() = %v, want 17", got)
	}
}
