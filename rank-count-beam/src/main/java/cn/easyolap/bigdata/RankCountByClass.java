package cn.easyolap.bigdata;

import cn.easyolap.bigdata.bean.Trip;
import cn.easyolap.bigdata.constant.RankType;
import cn.hutool.core.bean.BeanUtil;
import cn.hutool.core.date.DatePattern;
import cn.hutool.core.date.DateUtil;
import cn.hutool.core.util.RandomUtil;
import org.apache.beam.sdk.Pipeline;
import org.apache.beam.sdk.io.jdbc.JdbcIO;
import org.apache.beam.sdk.options.Default;
import org.apache.beam.sdk.options.Description;
import org.apache.beam.sdk.options.PipelineOptions;
import org.apache.beam.sdk.options.PipelineOptionsFactory;
import org.apache.beam.sdk.transforms.*;
import org.apache.beam.sdk.values.KV;
import org.apache.beam.sdk.values.PCollection;
import org.apache.beam.sdk.values.TypeDescriptor;
import org.apache.beam.sdk.values.TypeDescriptors;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigDecimal;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.time.LocalDate;
import java.util.Date;
import java.util.Iterator;
import java.util.List;
import java.util.UUID;


public class RankCountByClass {

    private static final Logger log = LoggerFactory.getLogger(RankCountByClass.class);




    public interface RankCountByClassOptions extends PipelineOptions {
        @Description("mysql conn url")
        @Default.String("jdbc:mysql://hostname:3306/test")
        String getConnUrl();
        void setConnUrl(String connUrl);

        @Description("mysql con username")
        @Default.String("root")
        String getUserName();
        void setUserName(String userName);

        @Description("mysql con password")
        @Default.String("")
        String getPassword();
        void setPassword(String password);

        @Description("top n setting")
        @Default.Integer(10)
        Integer getTopn();
        void setTopn(Integer topn);

        @Description("统计范围 开始日期，默认为当天")
        @Default.String("")
        String getScopeStart();
        void setScopeStart(String scopeStart);

        @Description("统计范围 结束日期，默认为当天")
        @Default.String("")
        String getScopeEnd();
        void setScopeEnd(String scopeEnd);

    }

    public static void runCount(RankCountByClassOptions options){
        Pipeline p = Pipeline.create(options);
        String driver = "com.mysql.cj.jdbc.Driver";
        String querySql = "select id, cid,ccid,sid sid,distance,run_time,avge_speed,max_speed,create_date " +
                "from u_trip t " ;
        String start = options.getScopeStart();
        String end  = options.getScopeEnd();
        if(start == null || start.length() == 0){
            start =  DateUtil.today();
            start += "00:00:01";
        }
        if(end == null || end.length() == 0){
            end = DateUtil.today();
            end += "59:59:59";
        }
        querySql +=   " where create_date >= '"+start +"' and create_date <= '"+end+"'";
        log.info("query sql:{}",querySql);
        //读取Mysql data
        PCollection<KV<String,Trip>> resultCollection =
                p.apply(JdbcIO.<KV<String,Trip>>read()
                .withDataSourceConfiguration(JdbcIO.DataSourceConfiguration.create(
                        driver,
                        options.getConnUrl())
                        .withUsername(options.getUserName())
                        .withPassword(options.getPassword()))
                .withQuery(querySql)
                // 对结果集中的每一条数据进行处理
                .withRowMapper(new JdbcIO.RowMapper<KV<String,Trip>>() {
                    @Override
                    public KV<String,Trip> mapRow(ResultSet resultSet) throws Exception {
                        Trip trip = new Trip();
                        String id = resultSet.getString("id");
                        String cid = resultSet.getString("cid");
                        String ccid = resultSet.getString("ccid");

                        String sid = resultSet.getString("sid");
                        BigDecimal distance = resultSet.getBigDecimal("distance");
                        Long runTime = resultSet.getLong("run_time");
                        Date createDate = resultSet.getDate("create_date");
                        if(createDate == null ){
                            createDate = new Date();
                        }
                        int year = DateUtil.year(createDate);
                        int month = DateUtil.month(createDate);
                        String day = DateUtil.format(createDate, DatePattern.NORM_DATE_PATTERN);
                        int week = DateUtil.weekOfYear(createDate);
                        trip.setId(id);
                        trip.setYear(year+"");
                        trip.setMonth(month+"");
                        trip.setDay(day);
                        trip.setWeek(week+"");

                        trip.setsId(sid);
                        trip.setCcId(ccid);
                        trip.setcId(cid);
                        trip.setDistance(distance);
                        trip.setRunTimes(runTime);
                        log.info("sid:{},distance:{}",sid , distance);
                        String key = sid +"_"+trip.getDay();
                        return KV.of(key,trip);
                    }
                }));

        // 根据sid聚合,把同一个人的数据聚合为一条
        PCollection<KV<String, Trip>> resultPerson = resultCollection.apply(
                GroupByKey.<String,Trip>create())
                // 对聚合后的结果进行处理
                .apply(MapElements.into(TypeDescriptors.kvs(TypeDescriptors.strings(),TypeDescriptor.of(Trip.class)))
                        .via(e -> {
                    Iterable<Trip> value = e.getValue();
                    if (value == null) {
                        throw new NullPointerException();
                    }
                    Iterator<Trip> iterator = e.getValue().iterator();
                    BigDecimal d = BigDecimal.ZERO;
                    long t = 0;
                    Trip tmp = null;
                    while (iterator.hasNext()) {
                        tmp = iterator.next();
                        d = d.add(tmp.getDistance());
                        t+= tmp.getRunTimes();
                    }
                    //把某个人的数据聚合完成
                   Trip trip = new Trip();
                    BeanUtil.copyProperties(tmp,trip);
                    //聚合后的
                    trip.setDistance(d);
                    trip.setRunTimes(t);
                    return KV.of(e.getKey(), trip);
                }));

                // 自定义算子打印结果集
        PCollection<KV<String, List<Trip>>> largest10ValuesPerKey =
                resultPerson.apply(ParDo.of(new DoFn<KV<String, Trip>, KV<String, Trip>>() {
                    @ProcessElement
                    public void processElement(ProcessContext context) {
                        // 从管道中取出的每个元素
                        KV<String, Trip> element = context.element();
                        log.info("========== trip element info:{}",element);
                        context.output(element);
                    }
                })).apply(
                        "Max top N",
                        Top.largestPerKey(options.getTopn()));


        PCollection<KV<String, Trip>> topnTrips = largest10ValuesPerKey.apply(
                ParDo.of(new DoFn<KV<String, List<Trip>>,KV<String, Trip>>() {
            @ProcessElement
            public void processElement(ProcessContext context) {
                // 从管道中取出的每个元素
                KV<String, List<Trip>> element = context.element();
                String key = element.getKey();
                List<Trip> vals = element.getValue();
                if(vals!=null){
                    log.info("===key is:{}==vals size is ：{}",key,vals.size());
                    int i = 1;
                    for(Trip t : vals){
                        Trip trip = new Trip();
                        trip.setRank(i);
                        BeanUtil.copyProperties(t,trip);
                        log.info("sid group by info:{},topn is:{}", key, i);

                        i++;
                        KV<String,  Trip> kv = KV.of(key,trip);
                        context.output(kv);
                    }

                }else{
                    log.info("=====vals is empty ");
                }
            }
        }));

        // 将结果集写入数据库
        topnTrips.apply(JdbcIO.<KV<String,Trip>>write()
                .withDataSourceConfiguration(JdbcIO.DataSourceConfiguration.create(
                        driver,
                        options.getConnUrl())
                        .withUsername(options.getUserName())
                        .withPassword(options.getPassword()))
                .withStatement("insert into u_ranks " +
                        "(year,month,day,week,sid,cid,ccid,distance,run_time,rank,types,id,create_date) " +
                        "values(?,?,?,?,?,?,?,?,?,?,?,?,?)")
                .withPreparedStatementSetter(new JdbcIO.PreparedStatementSetter<KV<String,Trip>>() {
                    @Override
                    public void setParameters(KV<String,Trip> e,
                                              PreparedStatement preparedStatement) throws Exception {
                        if (e != null) {
                            String key = e.getKey();
                            Trip element = e.getValue();

                            log.info("==========JdbcIO==============key:{},size:{}",key,element);
                            LocalDate createDate = element.getCreateDate();

                            preparedStatement.setString(1, element.getYear());
                            preparedStatement.setString(2, element.getMonth());
                            preparedStatement.setString(3, element.getDay());
                            preparedStatement.setString(4, element.getWeek());

                            preparedStatement.setString(5, element.getsId());
                            preparedStatement.setString(6, element.getcId());
                            preparedStatement.setString(7, element.getCcId());

                            preparedStatement.setBigDecimal(8, element.getDistance());
                            preparedStatement.setLong(9, element.getRunTimes());
                            preparedStatement.setLong(10, element.getRank());

                            preparedStatement.setInt(11, RankType.TYPE_CLASS_DAY);
                            preparedStatement.setString(12,key);
                            preparedStatement.setDate(13, new java.sql.Date(System.currentTimeMillis()));
                            //preparedStatement.execute();  不需要执行，如果加上此句则会出现数据重复

                        }
                    }
                }));

                p.run().waitUntilFinish();
    }



    public static void main(String[] args) {
        log.info("logback 成功了");
        RankCountByClass.RankCountByClassOptions options =
                PipelineOptionsFactory.fromArgs(args)
                        .withValidation()
                        .as(RankCountByClass.RankCountByClassOptions.class);

        runCount(options);
    }

}
