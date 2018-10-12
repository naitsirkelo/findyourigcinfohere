# README

The Project named FindYourIgcInfoHere will allow the user to access a deployed application on Heroku, 
and then gather required information regarding a logged flight from an IGC formatted file.

Deployed Heroku link:

https://igcinfoimt2681.herokuapp.com


**Main file: igc.go**

The file is commented all the way through in a minimalistic manner, to easiest explain the meaning 
behind the structure without disrupting the flow or making a clutter of the code.


**Libraries used:**

1. github.com/marni/goigc	
2. github.com/p3lim/iso8601	



The following sections of code are used as alternatives.


First section:

Alternative for formatting 'uptime' into ISO 8601,
if the included library would no longer be dependent.


```
S := int(time.Since(startTime) / time.Second)	// Gather number of seconds since start

Mi := S/60		// Calculate Minutes, Hours, Days, Months and Years based on seconds.
H :=  Mi/60
D :=  H/24
Mo := D/30
Y :=  Mo/12
              	                // Remove values when times reach max.			
for year := 0; year < Y; year ++ { 	       Mo -= 12  }
for month := 0; month < Mo; month ++ { 	   D  -= 30  }
for day := 0; day < D; day ++ { 	       H  -= 24  }
for hour := 0; hour < H; hour ++ { 	       Mi -= 60  }
for minute := 0; minute < Mi; minute ++ {  S  -= 60  } 

              	                // Create string with ISO 8601 format
time := "P"+strconv.Itoa(Y)+"Y"+strconv.Itoa(Mo)+"M"+strconv.Itoa(D)+"DT"+strconv.Itoa(H)+"H"+strconv.Itoa(Mi)+"M"+strconv.Itoa(S)+"S"

temp := MetaData{time, "Service for IGC tracks.", "v1"}
```

Second section:

Alternative version for ouput TrackIDs. This code encodes each struct instead
of printing the stored IDs as an int array.


```
for id := range TrackIds {
  temp := TrackId{id}
          // Encodes temporary struct and shows information on screen
  err := json.NewEncoder(w).Encode(temp)
  if(err != nil){
      fmt.Println("\nEncode Error:", err)
  }
}
```


Third section:

Internal IDs given to posted URLs are alternatively stored in an individual map TrackIds, instead of URL ID
being identification key in the TrackUrl map.
*This would instead be useful if the internal system did not use integeres for IDs.*
In case of ID as string the map could be set to 'map[int]string' and then linked with TrackUrl.

Downside to this approach:
Ineffective to store ID in this map when internal ID is an integer, so TrackUrl tracks its own key instead.

```
var TrackIds map[int]int

TrackIds = make(map[int]int)     // Initializing map arrays

TrackUrl[tempLen] = temp["url"].(string)
TrackIds[tempLen] = templen      // Storing URL ID
```
