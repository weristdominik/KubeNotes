package webhook

import (
	"encoding/json"
	"net/http"
	"fmt"
	"io"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ValidatePod (w http.ResponseWriter, r *http.Request) {
	var review admissionv1.AdmissionReview
	json.NewDecoder(r.Body).Decode(&review)

	pod := corev1.Pod{}
	json.Unmarshal(review.Request.Object.Raw, &pod)

	allowed := true
	message := ""

	if pod.Labels["zone"] == "" {
		allowed = false
		message = "Pod is missing 'zone' label"
	}

	review.Response = &admissionv1.AdmissionResponse{
		Allowed: allowed,
		UID:	 review.Request.UID,
		Result:	 &metav1.Status{
			Message: message,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)

	fmt.Println("Validated pod:", pod.Name, "Allowed:", allowed)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only GET allowed"))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func DebugRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---- DEBUG Admission Request ----")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// Note: use log for later better code
		fmt.Println("read body error:", err)
		return
	}
	defer r.Body.Close()

	// 1. Print raw JSON from kube-apiserver
	fmt.Println("RAW BODY:")
	fmt.Println(string(body))

	// 2. Decode AdmissionReview
	var review admissionv1.AdmissionReview
	if err := json.Unmarshal(body, &review); err != nil {
		// Note: use log for later better code
		fmt.Println("unmarshal admission review error:", err)
		return
	}

	// 3. Pretty-print the request
	b, _ := json.MarshalIndent(review.Request, "", "  ")
	fmt.Println("ADMISSION REQUEST:")
	fmt.Println(string(b))

	// 4. Extract the actual object (Pod, etc.)
	if review.Request != nil {
		fmt.Println("RAW OBJECT:")
		fmt.Println(string(review.Request.Object.Raw))
	}
}
