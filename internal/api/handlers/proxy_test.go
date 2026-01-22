package handlers

import "testing"

func TestJoinURLPathPreservesBase(t *testing.T) {
	cases := []struct {
		base string
		req  string
		want string
	}{
		{base: "/lab", req: "/api", want: "/lab/api"},
		{base: "/lab/", req: "api", want: "/lab/api"},
		{base: "", req: "/api", want: "/api"},
		{base: "/lab", req: "", want: "/lab"},
		{base: "", req: "", want: "/"},
	}

	for _, tc := range cases {
		got := joinURLPath(tc.base, tc.req)
		if got != tc.want {
			t.Fatalf("joinURLPath(%q, %q) = %q, want %q", tc.base, tc.req, got, tc.want)
		}
	}
}
