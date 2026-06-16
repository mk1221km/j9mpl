/* Generated from 'MetricsLoggerTest.nrx' 16 Jun 2026 23:48:06 [v5.10] */
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
  int logMetric_counterCount=0;
  java.lang.String logMetric_expectedEx=null;
  java.lang.Throwable caughtEx=null;
  java.lang.String thrownClass=null;
  int matched=0;
  java.lang.Class expectedClass=null;
  com.factory.metrics.FuzzInput getAverageMetric_val1=null;
  com.factory.metrics.FuzzInput getAverageMetric_val1_input=null;
  java.lang.String getAverageMetric_p1=null;
  com.factory.metrics.FuzzInput getAverageMetric_val2=null;
  com.factory.metrics.FuzzInput getAverageMetric_val2_input=null;
  java.lang.String getAverageMetric_p2=null;
  int getAverageMetric_counterCount=0;
  java.lang.String getAverageMetric_expectedEx=null;
  com.factory.metrics.FuzzInput initDatabase_val1=null;
  com.factory.metrics.FuzzInput initDatabase_val1_input=null;
  java.lang.String initDatabase_p1=null;
  int initDatabase_counterCount=0;
  java.lang.String initDatabase_expectedEx=null;
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
    logMetric_counterCount=0;
    if (logMetric_val1_input.isCounter!=0) 
     logMetric_counterCount++;
    if (logMetric_val2_input.isCounter!=0) 
     logMetric_counterCount++;
    if (logMetric_counterCount<=1) 
     {
      logMetric_expectedEx=(java.lang.String)null;
      if (logMetric_val1_input.isCounter!=0) 
       {
        logMetric_expectedEx=logMetric_val1_input.expected;
       }
      if (logMetric_val2_input.isCounter!=0) 
       {
        logMetric_expectedEx=logMetric_val2_input.expected;
       }
      {try{
       if (1==0) 
        com.factory.metrics.MetricsLoggerTest.dummySignal();
       if (logMetric_val1_input.val!=null) 
        logMetric_p1=logMetric_val1_input.val;
       logMetric_p2=logMetric_val2_input.rec;
       MetricsLogger.logMetric(logMetric_p1,logMetric_p2);
       if (logMetric_expectedEx!=null) 
        {
         netrexx.lang.RexxIO.Say("Assertion Failure in logMetric: counter-example bypassed validation (no exception thrown). Expected: "+logMetric_expectedEx);
         java.lang.System.exit(1);
        }
      }
      catch (java.lang.Throwable $16){caughtEx=$16;
       thrownClass=caughtEx.getClass().getName();
       if (logMetric_expectedEx==null) 
        {
         netrexx.lang.RexxIO.Say("Assertion Failure in logMetric: happy path regression (unexpected exception: "+thrownClass+": "+caughtEx.getMessage()+")");
         caughtEx.printStackTrace();
         java.lang.System.exit(1);
        }
       else 
        {
         matched=0;
         {try{
          expectedClass=java.lang.Class.forName(logMetric_expectedEx);
          if (expectedClass.isInstance((java.lang.Object)caughtEx)) 
           matched=1;
         }
         catch (java.lang.ClassNotFoundException $17){
          if (!thrownClass.equals(logMetric_expectedEx)) 
           matched=1;
         }}
         if (matched==0) 
          {
           netrexx.lang.RexxIO.Say("Assertion Failure in logMetric: caught "+thrownClass+" ("+caughtEx.getMessage()+") but expected "+logMetric_expectedEx);
           java.lang.System.exit(1);
          }
        }
      }}
     }
    }
   }/*logMetric_val2*/
   }
  }/*logMetric_val1*/
  netrexx.lang.RexxIO.Say("  Method logMetric boundary exhaustion completed.");
  netrexx.lang.RexxIO.Say("Testing method getAverageMetric...");
  {int $20=0;com.factory.metrics.FuzzInput[] $19=new com.factory.metrics.FuzzInput[dbPathBounds.length];synchronized(dbPathBounds){for(;;){if($20==$19.length)break;$19[$20]=dbPathBounds[dbPathBounds.length-1-$20];$20++;}}getAverageMetric_val1:for(;;){if(--$20<0)break;getAverageMetric_val1=(com.factory.metrics.FuzzInput)$19[$20];
   getAverageMetric_val1_input=getAverageMetric_val1;
   getAverageMetric_p1=(java.lang.String)null;
   {int $23=0;com.factory.metrics.FuzzInput[] $22=new com.factory.metrics.FuzzInput[stringBounds.length];synchronized(stringBounds){for(;;){if($23==$22.length)break;$22[$23]=stringBounds[stringBounds.length-1-$23];$23++;}}getAverageMetric_val2:for(;;){if(--$23<0)break;getAverageMetric_val2=(com.factory.metrics.FuzzInput)$22[$23];
    getAverageMetric_val2_input=getAverageMetric_val2;
    getAverageMetric_p2=(java.lang.String)null;
    getAverageMetric_counterCount=0;
    if (getAverageMetric_val1_input.isCounter!=0) 
     getAverageMetric_counterCount++;
    if (getAverageMetric_val2_input.isCounter!=0) 
     getAverageMetric_counterCount++;
    if (getAverageMetric_counterCount<=1) 
     {
      getAverageMetric_expectedEx=(java.lang.String)null;
      if (getAverageMetric_val1_input.isCounter!=0) 
       {
        getAverageMetric_expectedEx=getAverageMetric_val1_input.expected;
       }
      if (getAverageMetric_val2_input.isCounter!=0) 
       {
        getAverageMetric_expectedEx=getAverageMetric_val2_input.expected;
       }
      {try{
       if (1==0) 
        com.factory.metrics.MetricsLoggerTest.dummySignal();
       if (getAverageMetric_val1_input.val!=null) 
        getAverageMetric_p1=getAverageMetric_val1_input.val;
       if (getAverageMetric_val2_input.val!=null) 
        getAverageMetric_p2=getAverageMetric_val2_input.val;
       MetricsLogger.getAverageMetric(getAverageMetric_p1,getAverageMetric_p2);
       if (getAverageMetric_expectedEx!=null) 
        {
         netrexx.lang.RexxIO.Say("Assertion Failure in getAverageMetric: counter-example bypassed validation (no exception thrown). Expected: "+getAverageMetric_expectedEx);
         java.lang.System.exit(1);
        }
      }
      catch (java.lang.Throwable $24){caughtEx=$24;
       thrownClass=caughtEx.getClass().getName();
       if (getAverageMetric_expectedEx==null) 
        {
         netrexx.lang.RexxIO.Say("Assertion Failure in getAverageMetric: happy path regression (unexpected exception: "+thrownClass+": "+caughtEx.getMessage()+")");
         caughtEx.printStackTrace();
         java.lang.System.exit(1);
        }
       else 
        {
         matched=0;
         {try{
          expectedClass=java.lang.Class.forName(getAverageMetric_expectedEx);
          if (expectedClass.isInstance((java.lang.Object)caughtEx)) 
           matched=1;
         }
         catch (java.lang.ClassNotFoundException $25){
          if (!thrownClass.equals(getAverageMetric_expectedEx)) 
           matched=1;
         }}
         if (matched==0) 
          {
           netrexx.lang.RexxIO.Say("Assertion Failure in getAverageMetric: caught "+thrownClass+" ("+caughtEx.getMessage()+") but expected "+getAverageMetric_expectedEx);
           java.lang.System.exit(1);
          }
        }
      }}
     }
    }
   }/*getAverageMetric_val2*/
   }
  }/*getAverageMetric_val1*/
  netrexx.lang.RexxIO.Say("  Method getAverageMetric boundary exhaustion completed.");
  netrexx.lang.RexxIO.Say("Testing method initDatabase...");
  {int $28=0;com.factory.metrics.FuzzInput[] $27=new com.factory.metrics.FuzzInput[dbPathBounds.length];synchronized(dbPathBounds){for(;;){if($28==$27.length)break;$27[$28]=dbPathBounds[dbPathBounds.length-1-$28];$28++;}}initDatabase_val1:for(;;){if(--$28<0)break;initDatabase_val1=(com.factory.metrics.FuzzInput)$27[$28];
   initDatabase_val1_input=initDatabase_val1;
   initDatabase_p1=(java.lang.String)null;
   initDatabase_counterCount=0;
   if (initDatabase_val1_input.isCounter!=0) 
    initDatabase_counterCount++;
   if (initDatabase_counterCount<=1) 
    {
     initDatabase_expectedEx=(java.lang.String)null;
     if (initDatabase_val1_input.isCounter!=0) 
      {
       initDatabase_expectedEx=initDatabase_val1_input.expected;
      }
     {try{
      if (1==0) 
       com.factory.metrics.MetricsLoggerTest.dummySignal();
      if (initDatabase_val1_input.val!=null) 
       initDatabase_p1=initDatabase_val1_input.val;
      MetricsLogger.initDatabase(initDatabase_p1);
      if (initDatabase_expectedEx!=null) 
       {
        netrexx.lang.RexxIO.Say("Assertion Failure in initDatabase: counter-example bypassed validation (no exception thrown). Expected: "+initDatabase_expectedEx);
        java.lang.System.exit(1);
       }
     }
     catch (java.lang.Throwable $29){caughtEx=$29;
      thrownClass=caughtEx.getClass().getName();
      if (initDatabase_expectedEx==null) 
       {
        netrexx.lang.RexxIO.Say("Assertion Failure in initDatabase: happy path regression (unexpected exception: "+thrownClass+": "+caughtEx.getMessage()+")");
        caughtEx.printStackTrace();
        java.lang.System.exit(1);
       }
      else 
       {
        matched=0;
        {try{
         expectedClass=java.lang.Class.forName(initDatabase_expectedEx);
         if (expectedClass.isInstance((java.lang.Object)caughtEx)) 
          matched=1;
        }
        catch (java.lang.ClassNotFoundException $30){
         if (!thrownClass.equals(initDatabase_expectedEx)) 
          matched=1;
        }}
        if (matched==0) 
         {
          netrexx.lang.RexxIO.Say("Assertion Failure in initDatabase: caught "+thrownClass+" ("+caughtEx.getMessage()+") but expected "+initDatabase_expectedEx);
          java.lang.System.exit(1);
         }
       }
     }}
    }
   }
  }/*initDatabase_val1*/
  netrexx.lang.RexxIO.Say("  Method initDatabase boundary exhaustion completed.");
  netrexx.lang.RexxIO.Say("=== [Phase III] Boundary Input Exhaustion Test Completed successfully! ===");
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