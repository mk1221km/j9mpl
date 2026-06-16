/* Generated from 'MetricsLoggerTest.nrx' 16 Jun 2026 23:11:34 [v5.10] */
/* Options: Annotations Binary Decimal Format Implicituses Java Logo Replace Trace2 Verbose3 */
package com.factory.metrics;
import com.factory.metrics.MetricRecord;
import java.sql.SQLException;
import com.factory.metrics.MetricsLogger;


public class MetricsLoggerTest{
 private static final java.lang.String $0="MetricsLoggerTest.nrx";
 
 @SuppressWarnings("unchecked") 
 
 public static void main(java.lang.String args[]){
  com.factory.metrics.FuzzInput stringBounds[];
  com.factory.metrics.FuzzInput dbPathBounds[];
  com.factory.metrics.FuzzInput doubleBounds[];
  com.factory.metrics.FuzzInput rexxBounds[];
  java.util.ArrayList recordBounds;
  com.factory.metrics.FuzzInput tsVal=null;
  com.factory.metrics.FuzzInput nameVal=null;
  com.factory.metrics.FuzzInput valVal=null;
  com.factory.metrics.MetricRecord rec=null;
  int recIsCounter=0;
  java.lang.String recExpected=null;
  int counterCount=0;
  com.factory.metrics.FuzzInput logMetric_val1=null;
  com.factory.metrics.FuzzInput logMetric_val1_input=null;
  java.lang.String logMetric_p1=null;
  java.lang.Object logMetric_val2=null;
  com.factory.metrics.RecordFuzzInput logMetric_val2_input=null;
  com.factory.metrics.MetricRecord logMetric_p2=null;
  int logMetric_isCounter=0;
  java.util.ArrayList logMetric_expectedExs=null;
  java.lang.Throwable logMetric_ex=null;
  java.lang.Throwable logMetric_caught=null;
  com.factory.metrics.FuzzInput getAverageMetric_val1=null;
  com.factory.metrics.FuzzInput getAverageMetric_val1_input=null;
  java.lang.String getAverageMetric_p1=null;
  com.factory.metrics.FuzzInput getAverageMetric_val2=null;
  com.factory.metrics.FuzzInput getAverageMetric_val2_input=null;
  java.lang.String getAverageMetric_p2=null;
  int getAverageMetric_isCounter=0;
  java.util.ArrayList getAverageMetric_expectedExs=null;
  java.lang.Throwable getAverageMetric_ex=null;
  java.lang.Throwable getAverageMetric_caught=null;
  com.factory.metrics.FuzzInput initDatabase_val1=null;
  com.factory.metrics.FuzzInput initDatabase_val1_input=null;
  java.lang.String initDatabase_p1=null;
  int initDatabase_isCounter=0;
  java.util.ArrayList initDatabase_expectedExs=null;
  java.lang.Throwable initDatabase_ex=null;
  java.lang.Throwable initDatabase_caught=null;
  netrexx.lang.RexxIO.Say("=== [Phase III] Starting Boundary Input Exhaustion Test for MetricsLogger ===");
  stringBounds=new com.factory.metrics.FuzzInput[]{new com.factory.metrics.FuzzInput("normal_string_test",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("cpu_usage",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("standard_channel",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("\' OR \'1\'=\'1",1,"java.sql.SQLException"),new com.factory.metrics.FuzzInput("\'; DROP TABLE system_metrics; --",1,"java.sql.SQLException"),new com.factory.metrics.FuzzInput("\' UNION SELECT null, null, null --",1,"java.sql.SQLException"),new com.factory.metrics.FuzzInput("generated/metrics_test.db",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("metrics.db",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("routing.db",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("/etc/passwd",1,"java.io.IOException"),new com.factory.metrics.FuzzInput("../../../etc/passwd",1,"java.io.IOException"),new com.factory.metrics.FuzzInput("C:\\Windows\\win.ini",1,"java.io.IOException"),new com.factory.metrics.FuzzInput("normal_string",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("",1,"java.lang.IllegalArgumentException"),new com.factory.metrics.FuzzInput((java.lang.String)null,1,"java.lang.IllegalArgumentException"),new com.factory.metrics.FuzzInput("   ",1,"java.lang.IllegalArgumentException")};
  dbPathBounds=new com.factory.metrics.FuzzInput[]{new com.factory.metrics.FuzzInput("generated/metricslogger_test.db",0,(java.lang.String)null),new com.factory.metrics.FuzzInput(":memory:",0,(java.lang.String)null),new com.factory.metrics.FuzzInput((java.lang.String)null,1,"java.lang.IllegalArgumentException")};
  doubleBounds=new com.factory.metrics.FuzzInput[]{new com.factory.metrics.FuzzInput("0",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("1",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("-1",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("1.5",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("-45.2",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("100.0",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("1.7976931348623157e+308",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("-1.7976931348623157e+308",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("1.7976931348623159e+308",1,"java.lang.NumberFormatException"),new com.factory.metrics.FuzzInput("-1.7976931348623159e+308",1,"java.lang.NumberFormatException")};
  rexxBounds=new com.factory.metrics.FuzzInput[]{new com.factory.metrics.FuzzInput("normal",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("0",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("1",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("-1",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("100",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("-50",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("999",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("2147483647",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("-2147483648",0,(java.lang.String)null),new com.factory.metrics.FuzzInput("2147483648",1,"java.lang.NumberFormatException"),new com.factory.metrics.FuzzInput("-2147483649",1,"java.lang.NumberFormatException")};
  recordBounds=new java.util.ArrayList();
  recordBounds.add((java.lang.Object)(new com.factory.metrics.RecordFuzzInput((com.factory.metrics.MetricRecord)null,1,"java.lang.IllegalArgumentException")));
  {int $3=0;com.factory.metrics.FuzzInput[] $2=new com.factory.metrics.FuzzInput[stringBounds.length];synchronized(stringBounds){for(;;){if($3==$2.length)break;$2[$3]=stringBounds[stringBounds.length-1-$3];$3++;}}tsVal:for(;;){if(--$3<0)break;tsVal=(com.factory.metrics.FuzzInput)$2[$3];
   {int $6=0;com.factory.metrics.FuzzInput[] $5=new com.factory.metrics.FuzzInput[stringBounds.length];synchronized(stringBounds){for(;;){if($6==$5.length)break;$5[$6]=stringBounds[stringBounds.length-1-$6];$6++;}}nameVal:for(;;){if(--$6<0)break;nameVal=(com.factory.metrics.FuzzInput)$5[$6];
    {int $9=0;com.factory.metrics.FuzzInput[] $8=new com.factory.metrics.FuzzInput[doubleBounds.length];synchronized(doubleBounds){for(;;){if($9==$8.length)break;$8[$9]=doubleBounds[doubleBounds.length-1-$9];$9++;}}valVal:for(;;){if(--$9<0)break;valVal=(com.factory.metrics.FuzzInput)$8[$9];
     rec=new MetricRecord();
     if (tsVal.val!=null) 
      rec.timestamp=tsVal.val;
     if (nameVal.val!=null) 
      rec.metricName=nameVal.val;
     if (valVal.val!=null) 
      rec.metricValue=new netrexx.lang.Rexx(valVal.val);
     recIsCounter=(tsVal.isCounter|nameVal.isCounter)|valVal.isCounter;
     recExpected=(java.lang.String)null;
     if (tsVal.isCounter!=0) 
      {
       recExpected=tsVal.expected;
      }
     else 
      if (nameVal.isCounter!=0) 
       {
        recExpected=nameVal.expected;
       }
      else 
       if (valVal.isCounter!=0) 
        {
         recExpected=valVal.expected;
        }
     counterCount=(tsVal.isCounter+nameVal.isCounter)+valVal.isCounter;
     if ((counterCount==0)|(counterCount==1)) 
      {
       recordBounds.add((java.lang.Object)(new com.factory.metrics.RecordFuzzInput(rec,recIsCounter,recExpected)));
      }
     }
    }/*valVal*/
    }
   }/*nameVal*/
   }
  }/*tsVal*/
  if (((((stringBounds!=null)&(dbPathBounds!=null))&(doubleBounds!=null))&(rexxBounds!=null))&(recordBounds!=null)) 
   netrexx.lang.RexxIO.Say("ok");
  netrexx.lang.RexxIO.Say("Testing method logMetric...");
  {int $12=0;com.factory.metrics.FuzzInput[] $11=new com.factory.metrics.FuzzInput[dbPathBounds.length];synchronized(dbPathBounds){for(;;){if($12==$11.length)break;$11[$12]=dbPathBounds[dbPathBounds.length-1-$12];$12++;}}logMetric_val1:for(;;){if(--$12<0)break;logMetric_val1=(com.factory.metrics.FuzzInput)$11[$12];
   logMetric_val1_input=logMetric_val1;
   logMetric_p1=(java.lang.String)null;
   {int $15=0;java.lang.Object[] $14=new java.lang.Object[recordBounds.size()];synchronized(recordBounds){java.util.Iterator $13=recordBounds.iterator();for(;;){if($15==$14.length)break;$14[$15]=$13.next();$15++;}}logMetric_val2:for(;;){if(--$15<0)break;logMetric_val2=(java.lang.Object)$14[$15];
    logMetric_val2_input=(com.factory.metrics.RecordFuzzInput)logMetric_val2;
    logMetric_p2=(com.factory.metrics.MetricRecord)null;
    logMetric_isCounter=0;
    logMetric_expectedExs=new java.util.ArrayList();
    if (logMetric_val1_input.isCounter!=0) 
     {
      logMetric_isCounter=1;
      logMetric_expectedExs.add((java.lang.Object)logMetric_val1_input.expected);
     }
    if (logMetric_val2_input.isCounter!=0) 
     {
      logMetric_isCounter=1;
      logMetric_expectedExs.add((java.lang.Object)logMetric_val2_input.expected);
     }
    logMetric_ex=(java.lang.Throwable)null;
    {try{
     if (1==0) 
      com.factory.metrics.MetricsLoggerTest.dummySignal();
     if (logMetric_val1_input.val!=null) 
      logMetric_p1=logMetric_val1_input.val;
     logMetric_p2=logMetric_val2_input.rec;
     MetricsLogger.logMetric(logMetric_p1,logMetric_p2);
    }
    catch (java.lang.Throwable $16){logMetric_caught=$16;
     logMetric_ex=logMetric_caught;
    }}
    com.factory.metrics.MetricsLoggerTest.assertResult("logMetric",logMetric_isCounter,logMetric_expectedExs,logMetric_ex);
    }
   }/*logMetric_val2*/
   }
  }/*logMetric_val1*/
  netrexx.lang.RexxIO.Say("  Method logMetric boundary exhaustion completed.");
  netrexx.lang.RexxIO.Say("Testing method getAverageMetric...");
  {int $19=0;com.factory.metrics.FuzzInput[] $18=new com.factory.metrics.FuzzInput[dbPathBounds.length];synchronized(dbPathBounds){for(;;){if($19==$18.length)break;$18[$19]=dbPathBounds[dbPathBounds.length-1-$19];$19++;}}getAverageMetric_val1:for(;;){if(--$19<0)break;getAverageMetric_val1=(com.factory.metrics.FuzzInput)$18[$19];
   getAverageMetric_val1_input=getAverageMetric_val1;
   getAverageMetric_p1=(java.lang.String)null;
   {int $22=0;com.factory.metrics.FuzzInput[] $21=new com.factory.metrics.FuzzInput[stringBounds.length];synchronized(stringBounds){for(;;){if($22==$21.length)break;$21[$22]=stringBounds[stringBounds.length-1-$22];$22++;}}getAverageMetric_val2:for(;;){if(--$22<0)break;getAverageMetric_val2=(com.factory.metrics.FuzzInput)$21[$22];
    getAverageMetric_val2_input=getAverageMetric_val2;
    getAverageMetric_p2=(java.lang.String)null;
    getAverageMetric_isCounter=0;
    getAverageMetric_expectedExs=new java.util.ArrayList();
    if (getAverageMetric_val1_input.isCounter!=0) 
     {
      getAverageMetric_isCounter=1;
      getAverageMetric_expectedExs.add((java.lang.Object)getAverageMetric_val1_input.expected);
     }
    if (getAverageMetric_val2_input.isCounter!=0) 
     {
      getAverageMetric_isCounter=1;
      getAverageMetric_expectedExs.add((java.lang.Object)getAverageMetric_val2_input.expected);
     }
    getAverageMetric_ex=(java.lang.Throwable)null;
    {try{
     if (1==0) 
      com.factory.metrics.MetricsLoggerTest.dummySignal();
     if (getAverageMetric_val1_input.val!=null) 
      getAverageMetric_p1=getAverageMetric_val1_input.val;
     if (getAverageMetric_val2_input.val!=null) 
      getAverageMetric_p2=getAverageMetric_val2_input.val;
     MetricsLogger.getAverageMetric(getAverageMetric_p1,getAverageMetric_p2);
    }
    catch (java.lang.Throwable $23){getAverageMetric_caught=$23;
     getAverageMetric_ex=getAverageMetric_caught;
    }}
    com.factory.metrics.MetricsLoggerTest.assertResult("getAverageMetric",getAverageMetric_isCounter,getAverageMetric_expectedExs,getAverageMetric_ex);
    }
   }/*getAverageMetric_val2*/
   }
  }/*getAverageMetric_val1*/
  netrexx.lang.RexxIO.Say("  Method getAverageMetric boundary exhaustion completed.");
  netrexx.lang.RexxIO.Say("Testing method initDatabase...");
  {int $26=0;com.factory.metrics.FuzzInput[] $25=new com.factory.metrics.FuzzInput[dbPathBounds.length];synchronized(dbPathBounds){for(;;){if($26==$25.length)break;$25[$26]=dbPathBounds[dbPathBounds.length-1-$26];$26++;}}initDatabase_val1:for(;;){if(--$26<0)break;initDatabase_val1=(com.factory.metrics.FuzzInput)$25[$26];
   initDatabase_val1_input=initDatabase_val1;
   initDatabase_p1=(java.lang.String)null;
   initDatabase_isCounter=0;
   initDatabase_expectedExs=new java.util.ArrayList();
   if (initDatabase_val1_input.isCounter!=0) 
    {
     initDatabase_isCounter=1;
     initDatabase_expectedExs.add((java.lang.Object)initDatabase_val1_input.expected);
    }
   initDatabase_ex=(java.lang.Throwable)null;
   {try{
    if (1==0) 
     com.factory.metrics.MetricsLoggerTest.dummySignal();
    if (initDatabase_val1_input.val!=null) 
     initDatabase_p1=initDatabase_val1_input.val;
    MetricsLogger.initDatabase(initDatabase_p1);
   }
   catch (java.lang.Throwable $27){initDatabase_caught=$27;
    initDatabase_ex=initDatabase_caught;
   }}
   com.factory.metrics.MetricsLoggerTest.assertResult("initDatabase",initDatabase_isCounter,initDatabase_expectedExs,initDatabase_ex);
   }
  }/*initDatabase_val1*/
  netrexx.lang.RexxIO.Say("  Method initDatabase boundary exhaustion completed.");
  netrexx.lang.RexxIO.Say("=== [Phase III] Boundary Input Exhaustion Test Completed successfully! ===");
  return;}
 
 
 @SuppressWarnings("unchecked") 
 
 public static void assertResult(java.lang.String methodName,int isCounter,java.util.ArrayList expectedExs,java.lang.Throwable ex){
  java.lang.String thrownEx=null;
  int npeExpected=0;
  int i=0;
  int matched=0;
  java.lang.String expEx=null;
  java.lang.Class expectedClass=null;
  int idx=0;
  java.lang.String exp=null;
  if (isCounter!=0) 
   {
    if (ex==null) 
     {
      netrexx.lang.RexxIO.Say("Assertion failure in "+methodName+": counter-example bypassed validation (no exception thrown)");
      java.lang.System.exit(1);
     }
    thrownEx=ex.getClass().getName();
    if (thrownEx.equals("java.lang.NullPointerException")) 
     {
      npeExpected=0;
      {int $28=(expectedExs.size())-1;i=0;i:for(;i<=$28;i++){
       if ((((java.lang.String)(expectedExs.get(i)))).equals("java.lang.NullPointerException")) 
        npeExpected=1;
       }
      }/*i*/
      if (npeExpected==0) 
       {
        netrexx.lang.RexxIO.Say("Assertion failure in "+methodName+": ungraceful crash (NullPointerException)");
        ex.printStackTrace();
        java.lang.System.exit(1);
       }
     }
    matched=0;
    {int $29=(expectedExs.size())-1;i=0;i:for(;i<=$29;i++){
     expEx=(java.lang.String)(expectedExs.get(i));
     {try{
      expectedClass=java.lang.Class.forName(expEx);
      if (expectedClass.isInstance((java.lang.Object)ex)) 
       matched=1;
     }
     catch (java.lang.ClassNotFoundException $30){
      ;
     }}
     }
    }/*i*/
    if (matched==0) 
     {
      netrexx.lang.RexxIO.Say("Assertion failure in "+methodName+": caught "+thrownEx+" ("+ex.getMessage()+") but none of the expected exceptions matched.");
      netrexx.lang.RexxIO.Say("Expected exceptions for "+methodName+":");
      {int $31=(expectedExs.size())-1;idx=0;idx:for(;idx<=$31;idx++){
       exp=(java.lang.String)(expectedExs.get(idx));
       if (exp==null) 
        exp="null";
       netrexx.lang.RexxIO.Say("  - "+exp);
       }
      }/*idx*/
      java.lang.System.exit(1);
     }
   }
  else 
   {
    if (ex!=null) 
     {
      netrexx.lang.RexxIO.Say("Assertion failure in "+methodName+": happy path regression (unexpected exception "+ex.getClass().getName()+": "+ex.getMessage()+")");
      ex.printStackTrace();
      java.lang.System.exit(1);
     }
   }
  return;}
 
 
 @SuppressWarnings("unchecked") 
 
 private static void dummySignal() throws java.lang.Throwable{
  ;
  return;}
 
 
 private MetricsLoggerTest(){return;}
 }


class FuzzInput{
 private static final java.lang.String $0="MetricsLoggerTest.nrx";
 /* properties public */
 public java.lang.String val;
 public int isCounter;
 public java.lang.String expected;
 
 @SuppressWarnings("unchecked") 
 
 public FuzzInput(java.lang.String aVal,int aIsCounter,java.lang.String aExpected){super();
  this.val=aVal;
  this.isCounter=aIsCounter;
  this.expected=aExpected;
  return;}
 
 }


class RecordFuzzInput{
 private static final java.lang.String $0="MetricsLoggerTest.nrx";
 /* properties public */
 public com.factory.metrics.MetricRecord rec;
 public int isCounter;
 public java.lang.String expected;
 
 @SuppressWarnings("unchecked") 
 
 public RecordFuzzInput(com.factory.metrics.MetricRecord aRec,int aIsCounter,java.lang.String aExpected){super();
  this.rec=aRec;
  this.isCounter=aIsCounter;
  this.expected=aExpected;
  return;}
 
 }