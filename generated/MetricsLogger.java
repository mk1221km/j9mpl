/* Generated from 'MetricsLogger.nrx' 16 Jun 2026 17:22:05 [v5.10] */
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
    catch (java.sql.SQLException $3){ex=$3;
     netrexx.lang.RexxIO.Say("Database average query error: "+ex.getMessage());
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
     catch (java.sql.SQLException $4){
      ;
     }}
    }}
   }
  return avg;
  }
 
 
 @SuppressWarnings("unchecked") 
 
 public static void main(java.lang.String args[]){
  return;}
 
 
 private MetricsLogger(){return;}
 }