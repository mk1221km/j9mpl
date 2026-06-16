/* Generated from 'MetricsLogger.nrx' 16 Jun 2026 17:55:05 [v5.10] */
/* Options: Annotations Binary Decimal Format Implicituses Java Logo Replace Trace2 Verbose3 */
package com.factory.metrics;
import java.sql.Connection;
import java.sql.ResultSet;
import java.sql.Statement;
import java.sql.PreparedStatement;
import java.sql.SQLException;
import java.sql.DriverManager;


class MetricsLoggerDummy{
 private static final java.lang.String $0="MetricsLogger.nrx";
 
 public MetricsLoggerDummy(){return;}
 }


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
  java.sql.Connection conn;
  java.sql.Statement stmt;
  java.sql.SQLException ex=null;
  conn=(java.sql.Connection)null;
  stmt=(java.sql.Statement)null;
  if ((dbPath!=null)&netrexx.lang.Rexx.toRexx(dbPath).OpNotEq(null,$01)) 
   {
    {try{
     conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
     stmt=conn.createStatement();
     stmt.executeUpdate("CREATE TABLE IF NOT EXISTS system_metrics ("+"timestamp TEXT, "+"name TEXT, "+"value REAL)");
    }
    catch (java.sql.SQLException $1){ex=$1;
     netrexx.lang.RexxIO.Say("Database init error: "+ex.getMessage());
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
  java.sql.Connection conn;
  java.sql.PreparedStatement stmt;
  java.sql.SQLException ex=null;
  conn=(java.sql.Connection)null;
  stmt=(java.sql.PreparedStatement)null;
  if ((dbPath!=null)&netrexx.lang.Rexx.toRexx(dbPath).OpNotEq(null,$01)) 
   {
    {try{
     conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
     stmt=conn.prepareStatement("INSERT INTO system_metrics (timestamp, name, value) VALUES (?, ?, ?)");
     stmt.setString(1,record.timestamp);
     stmt.setString(2,record.metricName);
     stmt.setDouble(3,record.metricValue.todouble());
     stmt.executeUpdate();
    }
    catch (java.sql.SQLException $3){ex=$3;
     netrexx.lang.RexxIO.Say("Error inserting metric: "+ex.getMessage());
    }
    finally{
     {try{
      if (stmt!=null) 
       stmt.close();
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
  java.sql.Connection conn;
  java.sql.PreparedStatement stmt;
  java.sql.ResultSet rs;
  netrexx.lang.Rexx avg;
  java.sql.SQLException ex=null;
  conn=(java.sql.Connection)null;
  stmt=(java.sql.PreparedStatement)null;
  rs=(java.sql.ResultSet)null;
  avg=new netrexx.lang.Rexx(0);
  if ((dbPath!=null)&netrexx.lang.Rexx.toRexx(dbPath).OpNotEq(null,$01)) 
   {
    {try{
     conn=DriverManager.getConnection("jdbc:sqlite:"+dbPath);
     stmt=conn.prepareStatement("SELECT AVG(value) FROM system_metrics WHERE name = ?");
     stmt.setString(1,name);
     rs=stmt.executeQuery();
     if (rs.next()) 
      {
       avg=new netrexx.lang.Rexx(rs.getDouble(1));
      }
    }
    catch (java.sql.SQLException $5){ex=$5;
     netrexx.lang.RexxIO.Say("Error fetching average metric: "+ex.getMessage());
    }
    finally{
     {try{
      if (rs!=null) 
       rs.close();
      if (stmt!=null) 
       stmt.close();
      if (conn!=null) 
       conn.close();
     }
     catch (java.sql.SQLException $6){
      ;
     }}
    }}
   }
  return avg;
  }
 
 
 @SuppressWarnings("unchecked") 
 
 public static void main(java.lang.String args[]){
  java.lang.String dbPath;
  com.factory.metrics.MetricRecord rec1=null;
  com.factory.metrics.MetricRecord rec2=null;
  netrexx.lang.Rexx avg=null;
  dbPath=(java.lang.String)null;
  if (args.length>0) 
   dbPath=args[0];
  else 
   dbPath="metrics.db";
  if ((dbPath!=null)&netrexx.lang.Rexx.toRexx(dbPath).OpNotEq(null,$01)) 
   {
    com.factory.metrics.MetricsLogger.initDatabase(dbPath);
    rec1=new com.factory.metrics.MetricRecord();
    rec1.timestamp="2025-01-01T00:00:00";
    rec1.metricName="cpu_usage";
    rec1.metricValue=new netrexx.lang.Rexx("45.5");
    com.factory.metrics.MetricsLogger.logMetric(dbPath,rec1);
    rec2=new com.factory.metrics.MetricRecord();
    rec2.timestamp="2025-01-01T00:01:00";
    rec2.metricName="cpu_usage";
    rec2.metricValue=new netrexx.lang.Rexx("50.2");
    com.factory.metrics.MetricsLogger.logMetric(dbPath,rec2);
    avg=com.factory.metrics.MetricsLogger.getAverageMetric(dbPath,"cpu_usage");
    netrexx.lang.RexxIO.Say(netrexx.lang.Rexx.toRexx("Average CPU usage: ").OpCc(null,avg));
   }
  else 
   {
    netrexx.lang.RexxIO.Say("Invalid database path");
   }
  return;}
 
 
 private MetricsLogger(){return;}
 }