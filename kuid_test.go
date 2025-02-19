package kuid

import (
	"strings"
	"sync"
	"testing"
)

func TestNewKUID(t *testing.T) {
	// Generate multiple KUIDs to ensure they're unique
	seen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		kuid, err := NewKUID()
		if err != nil {
			t.Fatalf("Failed to generate KUID: %v", err)
		}

		// Check string length
		str := kuid.String()
		if len(str) != size*2 {
			t.Errorf("Expected length %d, got %d", size*2, len(str))
		}

		// Check uniqueness
		if seen[str] {
			t.Errorf("Duplicate KUID generated: %s", str)
		}
		seen[str] = true

		// Verify string only contains valid base62 characters
		for _, c := range str {
			if !strings.ContainsRune(base62Chars, c) {
				t.Errorf("Invalid character in KUID string: %c", c)
			}
		}
	}
}

func TestKnownValues(t *testing.T) {
	tests := []struct {
		name string
		uuid string
	}{
		{
			name: "Zero UUID",
			uuid: "00000000-0000-0000-0000-000000000000",
		},
		{
			name: "Max UUID",
			uuid: "ffffffff-ffff-ffff-ffff-ffffffffffff",
		},
		{
			name: "Random UUID 1",
			uuid: "123e4567-e89b-12d3-a456-426614174000",
		},
		{
			name: "Random UUID 2",
			uuid: "d9db5cf3-c755-4f76-8746-04120f2644c6",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create KUID from UUID
			kuid, err := FromUUID(tt.uuid)
			if err != nil {
				t.Fatalf("FromUUID() error = %v", err)
			}

			// Convert to string and back to KUID
			str := kuid.String()
			decoded, err := FromString(str)
			if err != nil {
				t.Fatalf("FromString() error = %v", err)
			}

			// Verify roundtrip
			if !kuid.Equal(decoded) {
				t.Errorf("Roundtrip failed: values don't match")
				t.Logf("Original MSB: %x, LSB: %x", kuid.msb, kuid.lsb)
				t.Logf("Decoded MSB: %x, LSB: %x", decoded.msb, decoded.lsb)
			}

			// Verify UUID roundtrip
			uuidStr := decoded.ToUUID()
			if strings.ToLower(strings.ReplaceAll(tt.uuid, "-", "")) !=
				strings.ToLower(strings.ReplaceAll(uuidStr, "-", "")) {
				t.Errorf("UUID roundtrip failed")
				t.Logf("Original: %s", tt.uuid)
				t.Logf("Roundtrip: %s", uuidStr)
			}
		})
	}
}

func TestFromString(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		wantErr bool
	}{
		{
			name:    "Valid KUID string",
			str:     strings.Repeat("A", size*2),
			wantErr: false,
		},
		{
			name:    "Invalid length",
			str:     "ABC",
			wantErr: true,
		},
		{
			name:    "Invalid characters",
			str:     strings.Repeat("!", size*2),
			wantErr: true,
		},
		{
			name:    "Empty string",
			str:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kuid, err := FromString(tt.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				str := kuid.String()
				if str != tt.str {
					t.Errorf("String roundtrip failed: got %v, want %v", str, tt.str)
				}
			}
		})
	}
}

func TestKUID_Equal(t *testing.T) {
	kuid1, _ := NewKUID()
	kuid2, _ := NewKUID()
	kuid1Copy, _ := FromString(kuid1.String())

	tests := []struct {
		name string
		k1   *KUID
		k2   *KUID
		want bool
	}{
		{
			name: "Same KUID",
			k1:   kuid1,
			k2:   kuid1,
			want: true,
		},
		{
			name: "Different KUIDs",
			k1:   kuid1,
			k2:   kuid2,
			want: false,
		},
		{
			name: "KUID and its copy",
			k1:   kuid1,
			k2:   kuid1Copy,
			want: true,
		},
		{
			name: "KUID and nil",
			k1:   kuid1,
			k2:   nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.k1.Equal(tt.k2); got != tt.want {
				t.Errorf("KUID.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvalidByteSlices(t *testing.T) {
	tests := []struct {
		name    string
		bytes   []byte
		wantErr bool
	}{
		{
			name:    "Nil bytes",
			bytes:   nil,
			wantErr: true,
		},
		{
			name:    "Empty bytes",
			bytes:   []byte{},
			wantErr: true,
		},
		{
			name:    "Too short",
			bytes:   make([]byte, 15),
			wantErr: true,
		},
		{
			name:    "Too long",
			bytes:   make([]byte, 17),
			wantErr: true,
		},
		{
			name:    "Exact length",
			bytes:   make([]byte, 16),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FromBytes(tt.bytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConcurrentAccess(t *testing.T) {
	const numGoroutines = 100
	const numIterations = 100

	var wg sync.WaitGroup
	results := make(chan string, numGoroutines*numIterations)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numIterations; j++ {
				kuid, err := NewKUID()
				if err != nil {
					t.Errorf("Failed to generate KUID: %v", err)
					return
				}
				results <- kuid.String()
			}
		}()
	}

	wg.Wait()
	close(results)

	// Check for duplicates
	seen := make(map[string]bool)
	for result := range results {
		if seen[result] {
			t.Errorf("Duplicate KUID generated: %s", result)
		}
		seen[result] = true
	}
}

func BenchmarkKUID(b *testing.B) {
	b.Run("Generate", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := NewKUID()
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	sample, _ := NewKUID()
	sampleStr := sample.String()
	sampleUUID := sample.ToUUID()

	b.Run("ToString", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = sample.String()
		}
	})

	b.Run("FromString", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := FromString(sampleStr)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("FromUUID", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := FromUUID(sampleUUID)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
