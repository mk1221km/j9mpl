/* Generated from 'MetricsLogger.nrx' 15 Jun 2026 00:38:37 [v5.10] */
/* Options: Annotations Binary Decimal Format Implicituses Java Logo Replace Sourcedir Trace2 Verbose3 */
package com.factory.metrics;
import java.sql.Connection;
import java.sql.ResultSet;
import java.sql.Statement;
import java.sql.PreparedStatement;
import java.sql.SQLException;
import java.sql.DriverManager;


class MetricRecord{
 private static final java.lang.String $0="MetricsLogger.nrx";
 /* properties public */
 public java.lang.String timestamp;
 public java.lang.String metricName;
 public netrexx.lang.Rexx metricValue;
 
 public MetricRecord(){return;}
 }


public class MetricsLogger{
 private static final netrexx.lang.Rexx $01=netrexx.lang.Rexx.toRexx("null");
 private static final java.lang.String $0="MetricsLogger.nrx";
 
 @SuppressWarnings("unchecked") 
 
 public static void initDatabase(java.lang.String dbPath){
  java.sql.Connection conn=null;
  java.sql.Statement stmt=null;
  java.lang.Exception ex=null;
  if ((dbPath!=null)&netrexx.lang.Rexx.toRexx(dbPath).OpNotEq(null,$01)) 
   {
    conn=(java.sql.Connection)null;
    stmt=(java.sql.Statement)null;
    {try{
     java.lang.Class.forName("org.sqlite.JDBC");
     conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
     stmt=conn.createStatement();
     stmt.executeUpdate("CREATE TABLE IF NOT EXISTS system_metrics (timestamp TEXT, name TEXT, value REAL);");
    }
    catch (java.lang.Exception $1){ex=$1;
     ex.printStackTrace();
    }
    finally{
     {try{
      if (stmt!=null) 
       stmt.close();
      if (conn!=null) 
       conn.close();
     }
     catch (java.sql.SQLException $2){
      ;
     }}
    }}
   }
  return;}
 
 
 @SuppressWarnings("unchecked") 
 
 public static void logMetric(java.lang.String dbPath,com.factory.metrics.MetricRecord record){
  java.sql.Connection conn=null;
  java.sql.PreparedStatement pstmt=null;
  java.lang.Exception ex=null;
  if (((dbPath!=null)&netrexx.lang.Rexx.toRexx(dbPath).OpNotEq(null,$01))&(record!=null)) 
   {
    conn=(java.sql.Connection)null;
    pstmt=(java.sql.PreparedStatement)null;
    {try{
     java.lang.Class.forName("org.sqlite.JDBC");
     conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
     pstmt=conn.prepareStatement("INSERT INTO system_metrics (timestamp, name, value) VALUES (?, ?, ?);");
     pstmt.setString(1,record.timestamp);
     pstmt.setString(2,record.metricName);
     if (record.metricValue!=null) 
      pstmt.setDouble(3,record.metricValue.todouble());
     else 
      pstmt.setDouble(3,0.0D);
     pstmt.executeUpdate();
    }
    catch (java.lang.Exception $3){ex=$3;
     ex.printStackTrace();
    }
    finally{
     {try{
      if (pstmt!=null) 
       pstmt.close();
      if (conn!=null) 
       conn.close();
     }
     catch (java.sql.SQLException $4){
      ;
     }}
    }}
   }
  return;}
 
 
 @SuppressWarnings("unchecked") 
 
 public static netrexx.lang.Rexx getAverageMetric(java.lang.String dbPath,java.lang.String name){
  netrexx.lang.Rexx avgVal;
  java.sql.Connection conn=null;
  java.sql.PreparedStatement pstmt=null;
  java.sql.ResultSet rs=null;
  java.lang.Exception ex=null;
  avgVal=new netrexx.lang.Rexx(0);
  if (((dbPath!=null)&netrexx.lang.Rexx.toRexx(dbPath).OpNotEq(null,$01))&(name!=null)) 
   {
    conn=(java.sql.Connection)null;
    pstmt=(java.sql.PreparedStatement)null;
    rs=(java.sql.ResultSet)null;
    {try{
     java.lang.Class.forName("org.sqlite.JDBC");
     conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
     pstmt=conn.prepareStatement("SELECT AVG(value) FROM system_metrics WHERE name = ?;");
     pstmt.setString(1,name);
     rs=pstmt.executeQuery();
     if (rs.next()) 
      {
       avgVal=new netrexx.lang.Rexx(rs.getDouble(1));
      }
    }
    catch (java.lang.Exception $5){ex=$5;
     ex.printStackTrace();
    }
    finally{
     {try{
      if (rs!=null) 
       rs.close();
      if (pstmt!=null) 
       pstmt.close();
      if (conn!=null) 
       conn.close();
     }
     catch (java.sql.SQLException $6){
      ;
     }}
    }}
   }
  return avgVal;
  }
 
 
 @SuppressWarnings("unchecked") 
 
 public static void main(java.lang.String args[]){
  java.lang.String dbFile;
  com.factory.metrics.MetricRecord rec1;
  com.factory.metrics.MetricRecord rec2;
  netrexx.lang.Rexx avg;
  dbFile="generated/metrics.db";
  initDatabase(dbFile);
  rec1=new com.factory.metrics.MetricRecord();
  rec1.timestamp="2026-06-14T19:30:00";
  rec1.metricName="cpu_usage";
  rec1.metricValue=new netrexx.lang.Rexx(15.5F);
  logMetric(dbFile,rec1);
  rec2=new com.factory.metrics.MetricRecord();
  rec2.timestamp="2026-06-14T19:31:00";
  rec2.metricName="cpu_usage";
  rec2.metricValue=new netrexx.lang.Rexx(24.5F);
  logMetric(dbFile,rec2);
  avg=getAverageMetric(dbFile,"cpu_usage");
  netrexx.lang.RexxIO.Say(netrexx.lang.Rexx.toRexx("Average CPU Usage: ").OpCc(null,avg));
  return;}
 
 
 private MetricsLogger(){return;}
 }