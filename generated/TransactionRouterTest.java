/* Generated from 'TransactionRouterTest.nrx' 17 Jun 2026 00:48:00 [v5.10] */
/* Options: Annotations Binary Decimal Format Implicituses Java Logo Replace Trace2 Verbose3 */
package com.factory.routing;
import com.factory.routing.TransactionRecord;
import java.sql.SQLException;
import com.factory.routing.TransactionRouter;


public class TransactionRouterTest{
 private static final java.lang.String $0="TransactionRouterTest.nrx";
 
 @SuppressWarnings("unchecked") 
 
 public static void main(java.lang.String args[]){
  com.factory.routing.FuzzInput stringBounds[];
  com.factory.routing.FuzzInput dbPathBounds[];
  com.factory.routing.FuzzInput doubleBounds[];
  com.factory.routing.FuzzInput rexxBounds[];
  java.util.ArrayList recordBounds;
  com.factory.routing.FuzzInput tsVal=null;
  com.factory.routing.FuzzInput nameVal=null;
  com.factory.routing.FuzzInput valVal=null;
  com.factory.routing.TransactionRecord rec=null;
  int recIsCounter=0;
  java.lang.String recExpected=null;
  int counterCount=0;
  com.factory.routing.FuzzInput routeTransaction_val1=null;
  com.factory.routing.FuzzInput routeTransaction_val1_input=null;
  java.lang.String routeTransaction_p1=null;
  java.lang.Object routeTransaction_val2=null;
  com.factory.routing.RecordFuzzInput routeTransaction_val2_input=null;
  com.factory.routing.TransactionRecord routeTransaction_p2=null;
  int routeTransaction_counterCount=0;
  java.lang.String routeTransaction_expectedEx=null;
  java.lang.Throwable caughtEx=null;
  java.lang.String thrownClass=null;
  int matched=0;
  java.lang.Class expectedClass=null;
  com.factory.routing.FuzzInput initRoutingTable_val1=null;
  com.factory.routing.FuzzInput initRoutingTable_val1_input=null;
  java.lang.String initRoutingTable_p1=null;
  int initRoutingTable_counterCount=0;
  java.lang.String initRoutingTable_expectedEx=null;
  com.factory.routing.FuzzInput getTransactionCount_val1=null;
  com.factory.routing.FuzzInput getTransactionCount_val1_input=null;
  java.lang.String getTransactionCount_p1=null;
  com.factory.routing.FuzzInput getTransactionCount_val2=null;
  com.factory.routing.FuzzInput getTransactionCount_val2_input=null;
  java.lang.String getTransactionCount_p2=null;
  int getTransactionCount_counterCount=0;
  java.lang.String getTransactionCount_expectedEx=null;
  netrexx.lang.RexxIO.Say("=== [Phase III] Starting Boundary Input Exhaustion Test for TransactionRouter ===");
  stringBounds=new com.factory.routing.FuzzInput[]{new com.factory.routing.FuzzInput("normal_string_test",0,(java.lang.String)null),new com.factory.routing.FuzzInput("cpu_usage",0,(java.lang.String)null),new com.factory.routing.FuzzInput("standard_channel",0,(java.lang.String)null),new com.factory.routing.FuzzInput("\' OR \'1\'=\'1",1,"java.sql.SQLException"),new com.factory.routing.FuzzInput("\'; DROP TABLE system_metrics; --",1,"java.sql.SQLException"),new com.factory.routing.FuzzInput("\' UNION SELECT null, null, null --",1,"java.sql.SQLException"),new com.factory.routing.FuzzInput("generated/metrics_test.db",0,(java.lang.String)null),new com.factory.routing.FuzzInput("metrics.db",0,(java.lang.String)null),new com.factory.routing.FuzzInput("routing.db",0,(java.lang.String)null),new com.factory.routing.FuzzInput("/etc/passwd",1,"java.io.IOException"),new com.factory.routing.FuzzInput("../../../etc/passwd",1,"java.io.IOException"),new com.factory.routing.FuzzInput("C:\\Windows\\win.ini",1,"java.io.IOException"),new com.factory.routing.FuzzInput("normal_string",0,(java.lang.String)null),new com.factory.routing.FuzzInput("",1,"java.lang.IllegalArgumentException"),new com.factory.routing.FuzzInput((java.lang.String)null,1,"java.lang.IllegalArgumentException"),new com.factory.routing.FuzzInput("   ",1,"java.lang.IllegalArgumentException")};
  dbPathBounds=new com.factory.routing.FuzzInput[]{new com.factory.routing.FuzzInput("generated/transactionrouter_test.db",0,(java.lang.String)null),new com.factory.routing.FuzzInput(":memory:",0,(java.lang.String)null),new com.factory.routing.FuzzInput((java.lang.String)null,1,"java.lang.IllegalArgumentException")};
  doubleBounds=new com.factory.routing.FuzzInput[]{new com.factory.routing.FuzzInput("0",0,(java.lang.String)null),new com.factory.routing.FuzzInput("1",0,(java.lang.String)null),new com.factory.routing.FuzzInput("-1",0,(java.lang.String)null),new com.factory.routing.FuzzInput("1.5",0,(java.lang.String)null),new com.factory.routing.FuzzInput("-45.2",0,(java.lang.String)null),new com.factory.routing.FuzzInput("100.0",0,(java.lang.String)null),new com.factory.routing.FuzzInput("1.7976931348623157e+308",0,(java.lang.String)null),new com.factory.routing.FuzzInput("-1.7976931348623157e+308",0,(java.lang.String)null),new com.factory.routing.FuzzInput("1.7976931348623159e+308",1,"java.lang.NumberFormatException"),new com.factory.routing.FuzzInput("-1.7976931348623159e+308",1,"java.lang.NumberFormatException")};
  rexxBounds=new com.factory.routing.FuzzInput[]{new com.factory.routing.FuzzInput("normal",0,(java.lang.String)null),new com.factory.routing.FuzzInput("",0,(java.lang.String)null),new com.factory.routing.FuzzInput("0",0,(java.lang.String)null),new com.factory.routing.FuzzInput("1",0,(java.lang.String)null),new com.factory.routing.FuzzInput("-1",0,(java.lang.String)null),new com.factory.routing.FuzzInput("100",0,(java.lang.String)null),new com.factory.routing.FuzzInput("-50",0,(java.lang.String)null),new com.factory.routing.FuzzInput("999",0,(java.lang.String)null),new com.factory.routing.FuzzInput("2147483647",0,(java.lang.String)null),new com.factory.routing.FuzzInput("-2147483648",0,(java.lang.String)null),new com.factory.routing.FuzzInput("2147483648",1,"java.lang.NumberFormatException"),new com.factory.routing.FuzzInput("-2147483649",1,"java.lang.NumberFormatException")};
  recordBounds=new java.util.ArrayList();
  recordBounds.add((java.lang.Object)(new com.factory.routing.RecordFuzzInput((com.factory.routing.TransactionRecord)null,1,"java.lang.IllegalArgumentException")));
  {int $3=0;com.factory.routing.FuzzInput[] $2=new com.factory.routing.FuzzInput[stringBounds.length];synchronized(stringBounds){for(;;){if($3==$2.length)break;$2[$3]=stringBounds[stringBounds.length-1-$3];$3++;}}tsVal:for(;;){if(--$3<0)break;tsVal=(com.factory.routing.FuzzInput)$2[$3];
   {int $6=0;com.factory.routing.FuzzInput[] $5=new com.factory.routing.FuzzInput[stringBounds.length];synchronized(stringBounds){for(;;){if($6==$5.length)break;$5[$6]=stringBounds[stringBounds.length-1-$6];$6++;}}nameVal:for(;;){if(--$6<0)break;nameVal=(com.factory.routing.FuzzInput)$5[$6];
    {int $9=0;com.factory.routing.FuzzInput[] $8=new com.factory.routing.FuzzInput[doubleBounds.length];synchronized(doubleBounds){for(;;){if($9==$8.length)break;$8[$9]=doubleBounds[doubleBounds.length-1-$9];$9++;}}valVal:for(;;){if(--$9<0)break;valVal=(com.factory.routing.FuzzInput)$8[$9];
     rec=new TransactionRecord();
     if (tsVal.val!=null) 
      rec.priority=tsVal.val;
     if (nameVal.val!=null) 
      rec.receiver=nameVal.val;
     if (tsVal.val!=null) 
      rec.sender=tsVal.val;
     if (valVal.val!=null) 
      rec.amount=new netrexx.lang.Rexx(valVal.val);
     if (nameVal.val!=null) 
      rec.txId=nameVal.val;
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
       recordBounds.add((java.lang.Object)(new com.factory.routing.RecordFuzzInput(rec,recIsCounter,recExpected)));
      }
     }
    }/*valVal*/
    }
   }/*nameVal*/
   }
  }/*tsVal*/
  if (((((stringBounds!=null)&(dbPathBounds!=null))&(doubleBounds!=null))&(rexxBounds!=null))&(recordBounds!=null)) 
   netrexx.lang.RexxIO.Say("ok");
  netrexx.lang.RexxIO.Say("Testing method routeTransaction...");
  {int $12=0;com.factory.routing.FuzzInput[] $11=new com.factory.routing.FuzzInput[dbPathBounds.length];synchronized(dbPathBounds){for(;;){if($12==$11.length)break;$11[$12]=dbPathBounds[dbPathBounds.length-1-$12];$12++;}}routeTransaction_val1:for(;;){if(--$12<0)break;routeTransaction_val1=(com.factory.routing.FuzzInput)$11[$12];
   routeTransaction_val1_input=routeTransaction_val1;
   routeTransaction_p1=(java.lang.String)null;
   {int $15=0;java.lang.Object[] $14=new java.lang.Object[recordBounds.size()];synchronized(recordBounds){java.util.Iterator $13=recordBounds.iterator();for(;;){if($15==$14.length)break;$14[$15]=$13.next();$15++;}}routeTransaction_val2:for(;;){if(--$15<0)break;routeTransaction_val2=(java.lang.Object)$14[$15];
    routeTransaction_val2_input=(com.factory.routing.RecordFuzzInput)routeTransaction_val2;
    routeTransaction_p2=(com.factory.routing.TransactionRecord)null;
    routeTransaction_counterCount=0;
    if (routeTransaction_val1_input.isCounter!=0) 
     routeTransaction_counterCount++;
    if (routeTransaction_val2_input.isCounter!=0) 
     routeTransaction_counterCount++;
    if (routeTransaction_counterCount<=1) 
     {
      routeTransaction_expectedEx=(java.lang.String)null;
      if (routeTransaction_val1_input.isCounter!=0) 
       {
        routeTransaction_expectedEx=routeTransaction_val1_input.expected;
       }
      if (routeTransaction_val2_input.isCounter!=0) 
       {
        routeTransaction_expectedEx=routeTransaction_val2_input.expected;
       }
      {try{
       if (1==0) 
        com.factory.routing.TransactionRouterTest.dummySignal();
       if (routeTransaction_val1_input.val!=null) 
        routeTransaction_p1=routeTransaction_val1_input.val;
       routeTransaction_p2=routeTransaction_val2_input.rec;
       TransactionRouter.routeTransaction(routeTransaction_p1,routeTransaction_p2);
       if (routeTransaction_expectedEx!=null) 
        {
         netrexx.lang.RexxIO.Say("Assertion Failure in routeTransaction: counter-example bypassed validation (no exception thrown). Expected: "+routeTransaction_expectedEx);
         java.lang.System.exit(1);
        }
      }
      catch (java.lang.Throwable $16){caughtEx=$16;
       thrownClass=caughtEx.getClass().getName();
       if (routeTransaction_expectedEx==null) 
        {
         netrexx.lang.RexxIO.Say("Assertion Failure in routeTransaction: happy path regression (unexpected exception: "+thrownClass+": "+caughtEx.getMessage()+")");
         caughtEx.printStackTrace();
         java.lang.System.exit(1);
        }
       else 
        {
         matched=0;
         {try{
          expectedClass=java.lang.Class.forName(routeTransaction_expectedEx);
          if (expectedClass.isInstance((java.lang.Object)caughtEx)) 
           matched=1;
         }
         catch (java.lang.ClassNotFoundException $17){
          if (!thrownClass.equals(routeTransaction_expectedEx)) 
           matched=1;
         }}
         if (matched==0) 
          {
           netrexx.lang.RexxIO.Say("Assertion Failure in routeTransaction: caught "+thrownClass+" ("+caughtEx.getMessage()+") but expected "+routeTransaction_expectedEx);
           java.lang.System.exit(1);
          }
        }
      }}
     }
    }
   }/*routeTransaction_val2*/
   }
  }/*routeTransaction_val1*/
  netrexx.lang.RexxIO.Say("  Method routeTransaction boundary exhaustion completed.");
  netrexx.lang.RexxIO.Say("Testing method initRoutingTable...");
  {int $20=0;com.factory.routing.FuzzInput[] $19=new com.factory.routing.FuzzInput[dbPathBounds.length];synchronized(dbPathBounds){for(;;){if($20==$19.length)break;$19[$20]=dbPathBounds[dbPathBounds.length-1-$20];$20++;}}initRoutingTable_val1:for(;;){if(--$20<0)break;initRoutingTable_val1=(com.factory.routing.FuzzInput)$19[$20];
   initRoutingTable_val1_input=initRoutingTable_val1;
   initRoutingTable_p1=(java.lang.String)null;
   initRoutingTable_counterCount=0;
   if (initRoutingTable_val1_input.isCounter!=0) 
    initRoutingTable_counterCount++;
   if (initRoutingTable_counterCount<=1) 
    {
     initRoutingTable_expectedEx=(java.lang.String)null;
     if (initRoutingTable_val1_input.isCounter!=0) 
      {
       initRoutingTable_expectedEx=initRoutingTable_val1_input.expected;
      }
     {try{
      if (1==0) 
       com.factory.routing.TransactionRouterTest.dummySignal();
      if (initRoutingTable_val1_input.val!=null) 
       initRoutingTable_p1=initRoutingTable_val1_input.val;
      TransactionRouter.initRoutingTable(initRoutingTable_p1);
      if (initRoutingTable_expectedEx!=null) 
       {
        netrexx.lang.RexxIO.Say("Assertion Failure in initRoutingTable: counter-example bypassed validation (no exception thrown). Expected: "+initRoutingTable_expectedEx);
        java.lang.System.exit(1);
       }
     }
     catch (java.lang.Throwable $21){caughtEx=$21;
      thrownClass=caughtEx.getClass().getName();
      if (initRoutingTable_expectedEx==null) 
       {
        netrexx.lang.RexxIO.Say("Assertion Failure in initRoutingTable: happy path regression (unexpected exception: "+thrownClass+": "+caughtEx.getMessage()+")");
        caughtEx.printStackTrace();
        java.lang.System.exit(1);
       }
      else 
       {
        matched=0;
        {try{
         expectedClass=java.lang.Class.forName(initRoutingTable_expectedEx);
         if (expectedClass.isInstance((java.lang.Object)caughtEx)) 
          matched=1;
        }
        catch (java.lang.ClassNotFoundException $22){
         if (!thrownClass.equals(initRoutingTable_expectedEx)) 
          matched=1;
        }}
        if (matched==0) 
         {
          netrexx.lang.RexxIO.Say("Assertion Failure in initRoutingTable: caught "+thrownClass+" ("+caughtEx.getMessage()+") but expected "+initRoutingTable_expectedEx);
          java.lang.System.exit(1);
         }
       }
     }}
    }
   }
  }/*initRoutingTable_val1*/
  netrexx.lang.RexxIO.Say("  Method initRoutingTable boundary exhaustion completed.");
  netrexx.lang.RexxIO.Say("Testing method getTransactionCount...");
  {int $25=0;com.factory.routing.FuzzInput[] $24=new com.factory.routing.FuzzInput[dbPathBounds.length];synchronized(dbPathBounds){for(;;){if($25==$24.length)break;$24[$25]=dbPathBounds[dbPathBounds.length-1-$25];$25++;}}getTransactionCount_val1:for(;;){if(--$25<0)break;getTransactionCount_val1=(com.factory.routing.FuzzInput)$24[$25];
   getTransactionCount_val1_input=getTransactionCount_val1;
   getTransactionCount_p1=(java.lang.String)null;
   {int $28=0;com.factory.routing.FuzzInput[] $27=new com.factory.routing.FuzzInput[stringBounds.length];synchronized(stringBounds){for(;;){if($28==$27.length)break;$27[$28]=stringBounds[stringBounds.length-1-$28];$28++;}}getTransactionCount_val2:for(;;){if(--$28<0)break;getTransactionCount_val2=(com.factory.routing.FuzzInput)$27[$28];
    getTransactionCount_val2_input=getTransactionCount_val2;
    getTransactionCount_p2=(java.lang.String)null;
    getTransactionCount_counterCount=0;
    if (getTransactionCount_val1_input.isCounter!=0) 
     getTransactionCount_counterCount++;
    if (getTransactionCount_val2_input.isCounter!=0) 
     getTransactionCount_counterCount++;
    if (getTransactionCount_counterCount<=1) 
     {
      getTransactionCount_expectedEx=(java.lang.String)null;
      if (getTransactionCount_val1_input.isCounter!=0) 
       {
        getTransactionCount_expectedEx=getTransactionCount_val1_input.expected;
       }
      if (getTransactionCount_val2_input.isCounter!=0) 
       {
        getTransactionCount_expectedEx=getTransactionCount_val2_input.expected;
       }
      {try{
       if (1==0) 
        com.factory.routing.TransactionRouterTest.dummySignal();
       if (getTransactionCount_val1_input.val!=null) 
        getTransactionCount_p1=getTransactionCount_val1_input.val;
       if (getTransactionCount_val2_input.val!=null) 
        getTransactionCount_p2=getTransactionCount_val2_input.val;
       TransactionRouter.getTransactionCount(getTransactionCount_p1,getTransactionCount_p2);
       if (getTransactionCount_expectedEx!=null) 
        {
         netrexx.lang.RexxIO.Say("Assertion Failure in getTransactionCount: counter-example bypassed validation (no exception thrown). Expected: "+getTransactionCount_expectedEx);
         java.lang.System.exit(1);
        }
      }
      catch (java.lang.Throwable $29){caughtEx=$29;
       thrownClass=caughtEx.getClass().getName();
       if (getTransactionCount_expectedEx==null) 
        {
         netrexx.lang.RexxIO.Say("Assertion Failure in getTransactionCount: happy path regression (unexpected exception: "+thrownClass+": "+caughtEx.getMessage()+")");
         caughtEx.printStackTrace();
         java.lang.System.exit(1);
        }
       else 
        {
         matched=0;
         {try{
          expectedClass=java.lang.Class.forName(getTransactionCount_expectedEx);
          if (expectedClass.isInstance((java.lang.Object)caughtEx)) 
           matched=1;
         }
         catch (java.lang.ClassNotFoundException $30){
          if (!thrownClass.equals(getTransactionCount_expectedEx)) 
           matched=1;
         }}
         if (matched==0) 
          {
           netrexx.lang.RexxIO.Say("Assertion Failure in getTransactionCount: caught "+thrownClass+" ("+caughtEx.getMessage()+") but expected "+getTransactionCount_expectedEx);
           java.lang.System.exit(1);
          }
        }
      }}
     }
    }
   }/*getTransactionCount_val2*/
   }
  }/*getTransactionCount_val1*/
  netrexx.lang.RexxIO.Say("  Method getTransactionCount boundary exhaustion completed.");
  netrexx.lang.RexxIO.Say("=== [Phase III] Boundary Input Exhaustion Test Completed successfully! ===");
  return;}
 
 
 @SuppressWarnings("unchecked") 
 
 private static void dummySignal() throws java.lang.Throwable{
  ;
  return;}
 
 
 private TransactionRouterTest(){return;}
 }


class FuzzInput{
 private static final java.lang.String $0="TransactionRouterTest.nrx";
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
 private static final java.lang.String $0="TransactionRouterTest.nrx";
 /* properties public */
 public com.factory.routing.TransactionRecord rec;
 public int isCounter;
 public java.lang.String expected;
 
 @SuppressWarnings("unchecked") 
 
 public RecordFuzzInput(com.factory.routing.TransactionRecord aRec,int aIsCounter,java.lang.String aExpected){super();
  this.rec=aRec;
  this.isCounter=aIsCounter;
  this.expected=aExpected;
  return;}
 
 }