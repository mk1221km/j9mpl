/* Generated from 'MetricsLoggerTest.nrx' 15 Jun 2026 00:39:08 [v5.10] */
/* Options: Annotations Binary Decimal Format Implicituses Java Logo Replace Sourcedir Trace2 Verbose3 */
package com.factory.metrics;
import java.sql.SQLException;


public class MetricsLoggerTest{
 private static final netrexx.lang.Rexx $01=netrexx.lang.Rexx.toRexx("null");
 private static final java.lang.String $0="MetricsLoggerTest.nrx";
 
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
  com.factory.metrics.MetricRecord rec=null;
  java.lang.String val1=null;
  java.lang.String logMetric_p1=null;
  java.lang.Object val2=null;
  com.factory.metrics.MetricRecord logMetric_p2=null;
  java.lang.RuntimeException ex=null;
  java.lang.String getAverageMetric_p1=null;
  java.lang.String getAverageMetric_p2=null;
  java.lang.String initDatabase_p1=null;
  netrexx.lang.RexxIO.Say("=== [Phase III] Starting Boundary Input Exhaustion Test for MetricsLogger ===");
  stringBounds=new java.lang.String[]{"","normal_string_test","\'; DROP TABLE system_metrics; --","null"};
  dbPathBounds=new java.lang.String[]{"generated/metrics_test.db",":memory:","null"};
  doubleBounds=new netrexx.lang.Rexx[]{new netrexx.lang.Rexx((byte)0),new netrexx.lang.Rexx((byte)1),new netrexx.lang.Rexx((byte)-1),new netrexx.lang.Rexx(999999999),new netrexx.lang.Rexx(1.79e+308D),new netrexx.lang.Rexx(-(1.79e+308D))};
  rexxBounds=new netrexx.lang.Rexx[]{new netrexx.lang.Rexx((byte)0),new netrexx.lang.Rexx((byte)1),new netrexx.lang.Rexx((byte)-1),new netrexx.lang.Rexx("normal"),new netrexx.lang.Rexx("")};
  recordBounds=new java.util.ArrayList();
  recordBounds.add((java.lang.Object)null);
  {int $3=0;java.lang.String[] $2=new java.lang.String[stringBounds.length];synchronized(stringBounds){for(;;){if($3==$2.length)break;$2[$3]=stringBounds[stringBounds.length-1-$3];$3++;}}tsVal:for(;;){if(--$3<0)break;tsVal=(java.lang.String)$2[$3];
   {int $6=0;java.lang.String[] $5=new java.lang.String[stringBounds.length];synchronized(stringBounds){for(;;){if($6==$5.length)break;$5[$6]=stringBounds[stringBounds.length-1-$6];$6++;}}nameVal:for(;;){if(--$6<0)break;nameVal=(java.lang.String)$5[$6];
    {int $9=0;netrexx.lang.Rexx[] $8=new netrexx.lang.Rexx[doubleBounds.length];synchronized(doubleBounds){for(;;){if($9==$8.length)break;$8[$9]=doubleBounds[doubleBounds.length-1-$9];$9++;}}valVal:for(;;){if(--$9<0)break;valVal=(netrexx.lang.Rexx)$8[$9];
     rec=new com.factory.metrics.MetricRecord();
     if (netrexx.lang.Rexx.toRexx(tsVal).OpNotEq(null,$01)) 
      rec.timestamp=tsVal;
     if (netrexx.lang.Rexx.toRexx(nameVal).OpNotEq(null,$01)) 
      rec.metricName=nameVal;
     rec.metricValue=valVal;
     recordBounds.add((java.lang.Object)rec);
     }
    }/*valVal*/
    }
   }/*nameVal*/
   }
  }/*tsVal*/
  netrexx.lang.RexxIO.Say("Testing method logMetric...");
  {int $12=0;java.lang.String[] $11=new java.lang.String[dbPathBounds.length];synchronized(dbPathBounds){for(;;){if($12==$11.length)break;$11[$12]=dbPathBounds[dbPathBounds.length-1-$12];$12++;}}val1:for(;;){if(--$12<0)break;val1=(java.lang.String)$11[$12];
   logMetric_p1=(java.lang.String)null;
   if (netrexx.lang.Rexx.toRexx(val1).OpNotEq(null,$01)) 
    logMetric_p1=val1;
   {int $15=0;java.lang.Object[] $14=new java.lang.Object[recordBounds.size()];synchronized(recordBounds){java.util.Iterator $13=recordBounds.iterator();for(;;){if($15==$14.length)break;$14[$15]=$13.next();$15++;}}val2:for(;;){if(--$15<0)break;val2=(java.lang.Object)$14[$15];
    logMetric_p2=(com.factory.metrics.MetricRecord)null;
    logMetric_p2=(com.factory.metrics.MetricRecord)val2;
    {try{
     com.factory.metrics.MetricsLogger.logMetric(logMetric_p1,logMetric_p2);
    }
    catch (java.lang.RuntimeException $16){ex=$16;
     ;
    }}
    }
   }/*val2*/
   }
  }/*val1*/
  netrexx.lang.RexxIO.Say("  Method logMetric boundary exhaustion completed.");
  netrexx.lang.RexxIO.Say("Testing method getAverageMetric...");
  {int $19=0;java.lang.String[] $18=new java.lang.String[dbPathBounds.length];synchronized(dbPathBounds){for(;;){if($19==$18.length)break;$18[$19]=dbPathBounds[dbPathBounds.length-1-$19];$19++;}}val1:for(;;){if(--$19<0)break;val1=(java.lang.String)$18[$19];
   getAverageMetric_p1=(java.lang.String)null;
   if (netrexx.lang.Rexx.toRexx(val1).OpNotEq(null,$01)) 
    getAverageMetric_p1=val1;
   {int $22=0;java.lang.String[] $21=new java.lang.String[stringBounds.length];synchronized(stringBounds){for(;;){if($22==$21.length)break;$21[$22]=stringBounds[stringBounds.length-1-$22];$22++;}}val2:for(;;){if(--$22<0)break;val2=(java.lang.Object)((java.lang.String)$21[$22]);
    getAverageMetric_p2=(java.lang.String)null;
    if (val2!="null") 
     getAverageMetric_p2=(java.lang.String)val2;
    {try{
     com.factory.metrics.MetricsLogger.getAverageMetric(getAverageMetric_p1,getAverageMetric_p2);
    }
    catch (java.lang.RuntimeException $23){ex=$23;
     ;
    }}
    }
   }/*val2*/
   }
  }/*val1*/
  netrexx.lang.RexxIO.Say("  Method getAverageMetric boundary exhaustion completed.");
  netrexx.lang.RexxIO.Say("Testing method initDatabase...");
  {int $26=0;java.lang.String[] $25=new java.lang.String[dbPathBounds.length];synchronized(dbPathBounds){for(;;){if($26==$25.length)break;$25[$26]=dbPathBounds[dbPathBounds.length-1-$26];$26++;}}val1:for(;;){if(--$26<0)break;val1=(java.lang.String)$25[$26];
   initDatabase_p1=(java.lang.String)null;
   if (netrexx.lang.Rexx.toRexx(val1).OpNotEq(null,$01)) 
    initDatabase_p1=val1;
   {try{
    com.factory.metrics.MetricsLogger.initDatabase(initDatabase_p1);
   }
   catch (java.lang.RuntimeException $27){ex=$27;
    ;
   }}
   }
  }/*val1*/
  netrexx.lang.RexxIO.Say("  Method initDatabase boundary exhaustion completed.");
  netrexx.lang.RexxIO.Say("=== [Phase III] Boundary Input Exhaustion Test Completed successfully! ===");
  return;}
 
 
 private MetricsLoggerTest(){return;}
 }