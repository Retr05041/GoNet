package main

import "testing"

func TestAddClient(t *testing.T) {
  srv := &server{}
  newClient := client{}
  
  got := srv.AddClient(newClient)
  want := 1
  
  if got != want {
    t.Error("AddClient Test Failed")
  }
}
