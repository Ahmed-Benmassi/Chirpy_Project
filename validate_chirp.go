package main

import (
	"encoding/json"
	"net/http"
	"strings"
)



func (cfg *apiConfig)handelervalidatechirp(w http.ResponseWriter, r *http.Request){   // handel the request to validate the chirp and return a response with the validation result in json format
	type Chirprequest struct{
		Body string `json:"body"`                                        
	}
	type response struct{
		CleanedBody string `json:"cleaned_body"`
	}

	bad_words:=[]string{"kerfuffle","sharbert","fornax"}

	var req Chirprequest

	decoder :=json.NewDecoder(r.Body)                                     // NewDecoder returns a new decoder that reads from r. The decoder introduces its own buffering and may read data from r beyond the JSON values requested.
	err:=decoder.Decode(&req)                                             // Decode reads the next JSON-encoded value from its input and stores it in the value pointed to by v.
	if err!=nil {
		respondWithError(w,http.StatusBadRequest, "Something went wrong")
		w.WriteHeader(500)
		return 
	}

	if len(req.Body)>400{
		respondWithError(w,http.StatusBadRequest, "Chirp is too long")
		w.WriteHeader(200)
		return
	}

	w.WriteHeader(200)

	cleanedText := cleanwords(req.Body,bad_words)

	respondWithJSON(w, http.StatusOK, map[string]string{
		"cleaned_body": cleanedText,
	})
}


func respondWithError(w http.ResponseWriter, code int, msg string){
	respondWithJSON(w , code, map[string]string{
		"error" :msg,
   })

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}){  // set the content type to json and write the status code and encode the payload as json and write it to the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)                                      // Encode writes the JSON encoding of v to the stream, followed by a newline character.

}

func cleanwords(body string, bad_words []string) string {   // clean the chirp from any bad words and return the cleaned text

	
	words:=strings.Split(body," ")
	for i, word := range words {
		lower:=strings.ToLower(word)
        for _, bad := range bad_words {
            if lower == bad {
                words[i] = "****"
                break
            }
			
        }
    }
	newtext:=strings.Join(words," ")
	return newtext
}

