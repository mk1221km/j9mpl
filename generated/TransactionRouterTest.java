/* Generated from 'TransactionRouterTest.nrx' 16 Jun 2026 17:39:37 [v5.10] */
/* Options: Annotations Binary Decimal Format Implicituses Java Logo Replace Trace2 Verbose3 */
package com.factory.routing;
import java.sql.SQLException;


public class TransactionRouterTest{
 private static final java.lang.String $0="TransactionRouterTest.nrx";
 
 @SuppressWarnings("unchecked") 
 
 public static void main(java.lang.String args[]){
  java.lang.String stringBounds[];
  java.lang.String dbPathBounds[];
  netrexx.lang.Rexx doubleBounds[];
  netrexx.lang.Rexx rexxBounds[];
  java.util.ArrayList recordBounds;
  java.lang.String tsVal=null;
  java.lang.String nameVal=null;
  netrexx.lang.Rexx valVal=null;
  com.factory.routing.TransactionRecord rec=null;
  java.lang.String val1=null;
  java.lang.String routeTransaction_p1=null;
  java.lang.Object val2=null;
  com.factory.routing.TransactionRecord routeTransaction_p2=null;
  java.lang.String initRoutingTable_p1=null;
  java.lang.String getTransactionCount_p1=null;
  java.lang.String getTransactionCount_p2=null;
  netrexx.lang.RexxIO.Say("=== [Phase III] Starting Boundary Input Exhaustion Test for TransactionRouter ===");
  stringBounds=new java.lang.String[]{"normal_string_test","\' OR \'1\'=\'1","\'; DROP TABLE system_metrics; --","\' UNION SELECT null, null, null --","/etc/passwd","../../../etc/passwd","C:\\Windows\\win.ini",(java.lang.String)null,"","   "};
  dbPathBounds=new java.lang.String[]{"generated/metrics_test.db",":memory:","null"};
  doubleBounds=new netrexx.lang.Rexx[]{new netrexx.lang.Rexx((byte)0),new netrexx.lang.Rexx((byte)1),new netrexx.lang.Rexx((byte)-1),new netrexx.lang.Rexx(new netrexx.lang.Rexx("1.7976931348623157e+308")),new netrexx.lang.Rexx((new netrexx.lang.Rexx("1.7976931348623157e+308")).OpMinus(null))};
  rexxBounds=new netrexx.lang.Rexx[]{new netrexx.lang.Rexx("normal"),new netrexx.lang.Rexx(""),new netrexx.lang.Rexx((byte)0),new netrexx.lang.Rexx((byte)1),new netrexx.lang.Rexx((byte)-1),new netrexx.lang.Rexx(2147483647),new netrexx.lang.Rexx(-2147483648L),new netrexx.lang.Rexx(2147483648L),new netrexx.lang.Rexx(-2147483649L)};
  recordBounds=new java.util.ArrayList();
  recordBounds.add((java.lang.Object)null);
  {int $3=0;java.lang.String[] $2=new java.lang.String[stringBounds.length];synchronized(stringBounds){for(;;){if($3==$2.length)break;$2[$3]=stringBounds[stringBounds.length-1-$3];$3++;}}tsVal:for(;;){if(--$3<0)break;tsVal=(java.lang.String)$2[$3];
   {int $6=0;java.lang.String[] $5=new java.lang.String[stringBounds.length];synchronized(stringBounds){for(;;){if($6==$5.length)break;$5[$6]=stringBounds[stringBounds.length-1-$6];$6++;}}nameVal:for(;;){if(--$6<0)break;nameVal=(java.lang.String)$5[$6];
    {int $9=0;netrexx.lang.Rexx[] $8=new netrexx.lang.Rexx[doubleBounds.length];synchronized(doubleBounds){for(;;){if($9==$8.length)break;$8[$9]=doubleBounds[doubleBounds.length-1-$9];$9++;}}valVal:for(;;){if(--$9<0)break;valVal=(netrexx.lang.Rexx)$8[$9];
     rec=new com.factory.routing.TransactionRecord();
     if (tsVal!=null) 
      if (!tsVal.equals("null")) 
       rec.priority=tsVal;
     if (nameVal!=null) 
      if (!nameVal.equals("null")) 
       rec.receiver=nameVal;
     if (tsVal!=null) 
      if (!tsVal.equals("null")) 
       rec.sender=tsVal;
     rec.amount=valVal;
     if (nameVal!=null) 
      if (!nameVal.equals("null")) 
       rec.txId=nameVal;
     recordBounds.add((java.lang.Object)rec);
     }
    }/*valVal*/
    }
   }/*nameVal*/
   }
  }/*tsVal*/
  if (((((stringBounds!=null)&(dbPathBounds!=null))&(doubleBounds!=null))&(rexxBounds!=null))&(recordBounds!=null)) 
   netrexx.lang.RexxIO.Say("ok");
  netrexx.lang.RexxIO.Say("Testing method routeTransaction...");
  {int $12=0;java.lang.String[] $11=new java.lang.String[dbPathBounds.length];synchronized(dbPathBounds){for(;;){if($12==$11.length)break;$11[$12]=dbPathBounds[dbPathBounds.length-1-$12];$12++;}}val1:for(;;){if(--$12<0)break;val1=(java.lang.String)$11[$12];
   routeTransaction_p1=(java.lang.String)null;
   if (val1!=null) 
    if (!val1.equals("null")) 
     routeTransaction_p1=val1;
   {int $15=0;java.lang.Object[] $14=new java.lang.Object[recordBounds.size()];synchronized(recordBounds){java.util.Iterator $13=recordBounds.iterator();for(;;){if($15==$14.length)break;$14[$15]=$13.next();$15++;}}val2:for(;;){if(--$15<0)break;val2=(java.lang.Object)$14[$15];
    routeTransaction_p2=(com.factory.routing.TransactionRecord)null;
    routeTransaction_p2=(com.factory.routing.TransactionRecord)val2;
    {try{
     com.factory.routing.TransactionRouter.routeTransaction(routeTransaction_p1,routeTransaction_p2);
    }
    catch (java.lang.RuntimeException $16){
     ;
    }}
    }
   }/*val2*/
   }
  }/*val1*/
  netrexx.lang.RexxIO.Say("  Method routeTransaction boundary exhaustion completed.");
  netrexx.lang.RexxIO.Say("Testing method initRoutingTable...");
  {int $19=0;java.lang.String[] $18=new java.lang.String[dbPathBounds.length];synchronized(dbPathBounds){for(;;){if($19==$18.length)break;$18[$19]=dbPathBounds[dbPathBounds.length-1-$19];$19++;}}val1:for(;;){if(--$19<0)break;val1=(java.lang.String)$18[$19];
   initRoutingTable_p1=(java.lang.String)null;
   if (val1!=null) 
    if (!val1.equals("null")) 
     initRoutingTable_p1=val1;
   {try{
    com.factory.routing.TransactionRouter.initRoutingTable(initRoutingTable_p1);
   }
   catch (java.lang.RuntimeException $20){
    ;
   }}
   }
  }/*val1*/
  netrexx.lang.RexxIO.Say("  Method initRoutingTable boundary exhaustion completed.");
  netrexx.lang.RexxIO.Say("Testing method getTransactionCount...");
  {int $23=0;java.lang.String[] $22=new java.lang.String[dbPathBounds.length];synchronized(dbPathBounds){for(;;){if($23==$22.length)break;$22[$23]=dbPathBounds[dbPathBounds.length-1-$23];$23++;}}val1:for(;;){if(--$23<0)break;val1=(java.lang.String)$22[$23];
   getTransactionCount_p1=(java.lang.String)null;
   if (val1!=null) 
    if (!val1.equals("null")) 
     getTransactionCount_p1=val1;
   {int $26=0;java.lang.String[] $25=new java.lang.String[stringBounds.length];synchronized(stringBounds){for(;;){if($26==$25.length)break;$25[$26]=stringBounds[stringBounds.length-1-$26];$26++;}}val2:for(;;){if(--$26<0)break;val2=(java.lang.Object)((java.lang.String)$25[$26]);
    getTransactionCount_p2=(java.lang.String)null;
    if (val2!=null) 
     if (val2!="null") 
      getTransactionCount_p2=(java.lang.String)val2;
    {try{
     com.factory.routing.TransactionRouter.getTransactionCount(getTransactionCount_p1,getTransactionCount_p2);
    }
    catch (java.lang.RuntimeException $27){
     ;
    }}
    }
   }/*val2*/
   }
  }/*val1*/
  netrexx.lang.RexxIO.Say("  Method getTransactionCount boundary exhaustion completed.");
  netrexx.lang.RexxIO.Say("=== [Phase III] Boundary Input Exhaustion Test Completed successfully! ===");
  return;}
 
 
 private TransactionRouterTest(){return;}
 }