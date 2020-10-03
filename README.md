This is a Go program that writes to a Google Drive Spreadsheet. We can select what to write to which cell. In a second step it also appends two rows to the document. 

Usage:
- Start with the instructions in step 1 & 2 on 
  https://developers.google.com/drive/api/v3/quickstart/go.
- Go to your Google Drive and create a new Spredsheet. Right click on the filename 
  and choose 'Get link' This looks like this:

https://docs.google.com/spreadsheets/d/1qBE9NfyyyyyyyyyQkTy-nnnnnnn0/edit?usp=sharing

- Rename 'copy.json.example' to 'copy.json'.
- Copy the id, the part between '/spreadsheets/d/' and '/edit', a row like  
  '1qBE9NfyyyyyyyyyQkTy-nnnnnnn0'. Insert this as 'spreadsheetId' in conf.json.
- Run with 'go run WriteToGoogleDriveSpreadsheet.go'. If everything works your Spreadsheet will 
  get some new data. You can adjust what to write (myval := []interface{}{"OK!"}) and in which cell (writeRange := "F1").

