/* Generated from 'MetricsLogger.nrx' 14 Jun 2026 20:32:46 [v5.10] */
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
 private static final java.lang.String $0="MetricsLogger.nrx";
 
 @SuppressWarnings("unchecked") 
 
 public static void initDatabase(java.lang.String dbPath){
  java.sql.Connection conn;
  java.sql.Statement stmt;
  java.sql.SQLException ex=null;
  conn=(java.sql.Connection)null;
  stmt=(java.sql.Statement)null;
  {try{
   conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
   stmt=conn.createStatement();
   stmt.executeUpdate("CREATE TABLE IF NOT EXISTS system_metrics (timestamp TEXT, name TEXT, value REAL);");
  }
  catch (java.sql.SQLException $1){ex=$1;
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
  return;}
 
 
 @SuppressWarnings("unchecked") 
 
 public static void logMetric(java.lang.String dbPath,com.factory.metrics.MetricRecord rec){
  java.sql.Connection conn;
  java.sql.PreparedStatement pstmt;
  java.sql.SQLException ex=null;
  conn=(java.sql.Connection)null;
  pstmt=(java.sql.PreparedStatement)null;
  {try{
   conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
   pstmt=conn.prepareStatement("INSERT INTO system_metrics (timestamp, name, value) VALUES (?, ?, ?);");
   pstmt.setString(1,rec.timestamp);
   pstmt.setString(2,rec.metricName);
   pstmt.setDouble(3,rec.metricValue.todouble());
   pstmt.executeUpdate();
  }
  catch (java.sql.SQLException $3){ex=$3;
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
  return;}
 
 
 @SuppressWarnings("unchecked") 
 
 public static netrexx.lang.Rexx getAverageMetric(java.lang.String dbPath,java.lang.String name){
  java.sql.Connection conn;
  java.sql.PreparedStatement pstmt;
  java.sql.ResultSet rs;
  netrexx.lang.Rexx avgVal;
  java.sql.SQLException ex=null;
  conn=(java.sql.Connection)null;
  pstmt=(java.sql.PreparedStatement)null;
  rs=(java.sql.ResultSet)null;
  avgVal=new netrexx.lang.Rexx(0);
  {try{
   conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
   pstmt=conn.prepareStatement("SELECT AVG(value) FROM system_metrics WHERE name = ?;");
   pstmt.setString(1,name);
   rs=pstmt.executeQuery();
   if (rs.next()) 
    {
     avgVal=new netrexx.lang.Rexx(rs.getDouble(1));
    }
  }
  catch (java.sql.SQLException $5){ex=$5;
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