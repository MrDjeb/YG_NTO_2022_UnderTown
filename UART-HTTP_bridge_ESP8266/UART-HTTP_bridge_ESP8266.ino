#include <EEPROM.h>
#include <ArduinoJson.h>
#include <SPI.h>
//#include <MFRC522.h>
#include <ESP8266HTTPClient.h>
#include <ESP8266WiFi.h>

//++++++++++НАСТРОЙКИ++++++++++//
#define MASTER_WRITE false
#define MASTER_HASH "1d6e126e3707e5ed51c5eea083e82bf4"

#define YOUR_SSID "9d614f1d7b348fca48e39cb91db399e0"           
#define YOUR_PASS "25d55ad283aa400af464c76d713c07ad"                      
#define SERVER_LOCATION "http://192.168.137.1:39999/"        

#define CONTENT_TYPE "application/json"           //MEMI-тип запроса text/plain
#define BLOCK1 60                                 //Диапозон блоков с которыми работаем, ниж. граница
#define BLOCK2 61                                 //Диапозон блоков с которыми работаем, верх. граница
#define SECRET "777888999"                        //Пароль доступа куда шлём

#define RST_PIN 5      //Пин куда подключаем RST
#define SS_PIN 15      //Пин куда подключаем SDA
#define E_ADDRESS 32   //Адрес ячейки в EPROOM, где записан ID esp
//++++++++++НАСТРОЙКИ++++++++++//

boolean Master = false;
String MasterHash = "nill";
String new_hash;
String IDesp;

#define ASCII_CONVERT 0
byte buff[64];

StaticJsonDocument<200> Make_Json(String StrSendJ){
  byte len = StrSendJ[0];
  StrSendJ = StrSendJ.substring(StrSendJ.indexOf('{')+1, StrSendJ.indexOf('}'));
  String Delimeters[len];
  int i = 0;
  char Delby = ':';
  while (StrSendJ.indexOf(Delby) >= 0) {
   int delim = StrSendJ.indexOf(Delby);
   Delimeters[i] = (StrSendJ.substring(0, delim));
   StrSendJ = StrSendJ.substring(delim + 1, StrSendJ.length());
   i++;
   if (StrSendJ.indexOf(Delby) == -1) {
     Delimeters[i] = StrSendJ;
   }
  }
  StaticJsonDocument<200> SendJ;
  for(int i=0;i<len;i+=2){
    SendJ[Delimeters[i]]=Delimeters[i+1];
  }
  for(int i=0;i<len;i++){
    Serial.println(Delimeters[i]);
  }
  return SendJ;
}

StaticJsonDocument<200> Post_Json(StaticJsonDocument<200> SendJ){
  Serial.print("Finished json line: "); serializeJson(SendJ, Serial); Serial.println();
  String StrSendJ;
  serializeJson(SendJ, StrSendJ);

  WiFiClient client;
  HTTPClient http;
  http.begin(client, String(SERVER_LOCATION)+"esp");
  http.addHeader("Content-Type", CONTENT_TYPE);
  int httpCode = http.POST(StrSendJ);
  String StrAcceptJ = http.getString();
  Serial.print("Accepted json line: "); Serial.println(StrAcceptJ);
  if (httpCode != 200) {Serial.print("Http code: "); Serial.println(httpCode);}
  http.end();
    
  StaticJsonDocument<200> AcceptJ;
  DeserializationError error = deserializeJson(AcceptJ, StrAcceptJ);
  if (error) {Serial.print("DeserializeJson failed: "); Serial.println(error.c_str());}
  return AcceptJ;
}



String strData = "";
boolean recievedFlag;

void setup() {
  Serial.begin(115200);
  pinMode(13,OUTPUT);
  digitalWrite(13,LOW);
  WiFi.begin(YOUR_SSID, YOUR_PASS);
  
  //++++++++++ПОДКЛЮЧЕНИЕ К WiFi++++++++++//
  Serial.println(); Serial.print("Waiting for connection");
  while (WiFi.status() != WL_CONNECTED) {
    Serial.print(".");
    delay(50);
  }
  Serial.println(); Serial.println("Wifi is connect!");
  //++++++++++ПОДКЛЮЧЕНИЕ К WiFi++++++++++//
}



void loop() {
  if (WiFi.status() == WL_CONNECTED) {
    String UART=getUART();
    if (UART!=""){
      Serial.println(UART);
      Post_Json(Make_Json(UART));
    }
  } else {
    Serial.println("Error in WiFi connection!");
    delay(500);
  }

}
String getUART(){
  while (Serial.available()) {        
    //digitalWrite(13,HIGH);    
    strData += (char)Serial.read();        
    recievedFlag = true;                   
    delay(5);                           
  }
  //digitalWrite(13,LOW); 
  if (recievedFlag) {
    String ans="";                      
    ans = strData;               
    strData = "";                          
    recievedFlag = false; 
    return ans;                
  }
  return "";
}
