/* Generated from 'TelemetryEngine.nrx' 14 Jun 2026 18:29:48 [v5.10] */
/* Options: Annotations Decimal Format Implicituses Java Logo Replace Sourcedir Trace2 Verbose3 */
package com.factory.telemetry;
import com.sun.net.httpserver.HttpExchange;


class TelemetryRecord{
 private static final java.lang.String $0="TelemetryEngine.nrx";
 /* properties public */
 public java.lang.String deviceId;
 public netrexx.lang.Rexx voltage;
 public java.lang.String status;
 
 public TelemetryRecord(){return;}
 }


public class TelemetryEngine{
 private static final netrexx.lang.Rexx $01=new netrexx.lang.Rexx("12.0");
 private static final java.lang.String $0="TelemetryEngine.nrx";
 
 @SuppressWarnings("unchecked") 
 
 public static void processData(com.sun.net.httpserver.HttpExchange exchange,com.factory.telemetry.TelemetryRecord rec){
  if (rec.voltage.OpLt(null,$01)) 
   {
    rec.status="FAULT: UNDERVOLTAGE";
   }
  else 
   {
    rec.status="OPERATIONAL";
   }
  return;}
 
 
 private TelemetryEngine(){return;}
 }