package xchango

type exchange2006 struct{}

func (exchange2006) FolderRequest() string {
	return `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:typ="http://schemas.microsoft.com/exchange/services/2006/types" xmlns:mes="http://schemas.microsoft.com/exchange/services/2006/messages">
<soapenv:Header>
</soapenv:Header>
<soapenv:Body>
<mes:GetFolder>
<mes:FolderShape>
<typ:BaseShape>IdOnly</typ:BaseShape>
</mes:FolderShape>
<mes:FolderIds>
<typ:DistinguishedFolderId Id="calendar" />
</mes:FolderIds>
</mes:GetFolder>
</soapenv:Body>
</soapenv:Envelope>`
}

func (exchange2006) CalendarRequest() string {
	return `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:typ="http://schemas.microsoft.com/exchange/services/2006/types" xmlns:mes="http://schemas.microsoft.com/exchange/services/2006/messages">
   <soapenv:Header>
   </soapenv:Header>
   <soapenv:Body>
      <mes:FindItem Traversal="Shallow">
         <mes:ItemShape>
            <typ:BaseShape>IdOnly</typ:BaseShape>
         </mes:ItemShape>
         <!--You have a CHOICE of the next 4 items at this level-->
         <mes:CalendarView MaxEntriesReturned="{{ .MaxFetchSize }}" StartDate="{{ .StartDate }}" EndDate="{{ .EndDate }}"/>
         <mes:ParentFolderIds>
           <!--You have a CHOICE of the next 2 items at this level-->
            <typ:FolderId Id="{{ .FolderId }}" ChangeKey="{{ .ChangeKey }}" />
         </mes:ParentFolderIds>
      </mes:FindItem>
   </soapenv:Body>
</soapenv:Envelope>`
}

func (exchange2006) CalendarDetailRequest() string {
	return `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:typ="http://schemas.microsoft.com/exchange/services/2006/types" xmlns:mes="http://schemas.microsoft.com/exchange/services/2006/messages">
<soapenv:Header>
</soapenv:Header>
<soapenv:Body>
<mes:GetItem>
    <mes:ItemShape>
        <typ:BaseShape>IdOnly</typ:BaseShape>
        <typ:AdditionalProperties>
          <typ:FieldURI FieldURI="item:Subject" />
          <typ:FieldURI FieldURI="calendar:Start" />
          <typ:FieldURI FieldURI="calendar:End" />
          <typ:FieldURI FieldURI="item:DisplayTo" />
          <typ:FieldURI FieldURI="item:DisplayCc" />
          <typ:FieldURI FieldURI="calendar:IsAllDayEvent" />
          <typ:FieldURI FieldURI="calendar:Location" />
          <typ:FieldURI FieldURI="calendar:MyResponseType" />
          <typ:FieldURI FieldURI="calendar:Organizer" />
          <typ:FieldURI FieldURI="item:Body" />
        </typ:AdditionalProperties>
    </mes:ItemShape>
    <mes:ItemIds>
      {{ range .Appointments }}<typ:ItemId Id="{{ .ItemId }}" ChangeKey="{{ .ChangeKey}}" />
    {{ end }}</mes:ItemIds>
</mes:GetItem>
</soapenv:Body>
</soapenv:Envelope>`
}
