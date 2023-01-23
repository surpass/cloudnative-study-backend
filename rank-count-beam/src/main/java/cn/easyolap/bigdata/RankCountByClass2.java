package cn.easyolap.bigdata;

import org.apache.beam.sdk.Pipeline;
import org.apache.beam.sdk.io.jdbc.JdbcIO;
import org.apache.beam.sdk.options.Default;
import org.apache.beam.sdk.options.Description;
import org.apache.beam.sdk.options.PipelineOptions;
import org.apache.beam.sdk.options.PipelineOptionsFactory;
import org.apache.beam.sdk.transforms.DoFn;
import org.apache.beam.sdk.transforms.GroupByKey;
import org.apache.beam.sdk.transforms.MapElements;
import org.apache.beam.sdk.transforms.ParDo;
import org.apache.beam.sdk.values.KV;
import org.apache.beam.sdk.values.PCollection;
import org.apache.beam.sdk.values.TypeDescriptors;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.util.Iterator;


public class RankCountByClass2 {

    private static final Logger log = LoggerFactory.getLogger(RankCountByClass2.class);

    public interface RankCountByClassOptions extends PipelineOptions {
        @Description("mysql con url")
        @Default.String("jdbc:mysql://hostname:3306/test")
        String getConnUrl();

        String setConnUrl(String url);

        @Description("mysql con username")
        @Default.String("root")
        String getUserName();
        String setUserName(String userName);

        @Description("mysql con password")
        @Default.String("")
        String getPassword();
        String setPassword(String password);


    }

    public static void runCount(RankCountByClassOptions options){
        Pipeline p = Pipeline.create(options);

        //读取Mysql
        PCollection<KV<String, Integer>> resultCollection =
                p.apply(JdbcIO.<KV<String, Integer>>read()
                .withDataSourceConfiguration(JdbcIO.DataSourceConfiguration.create(
                        "com.mysql.jdbc.Driver", options.getConnUrl())
                        .withUsername(options.getUserName())
                        .withPassword(options.getPassword()))
                .withQuery("select * from u_trip")
                // 对结果集中的每一条数据进行处理
                .withRowMapper(new JdbcIO.RowMapper<KV<String, Integer>>() {
                    @Override
                    public KV<String, Integer> mapRow(ResultSet resultSet) throws Exception {
                        String id = resultSet.getString(1);
                        String name = resultSet.getString(2);
                        log.info("id:{},name:{}",id , name);
                        return KV.of(name, 1);
                    }
                }));

                // 根据key聚合
        PCollection<KV<String, Integer>> result = resultCollection.apply(
                GroupByKey.<String,Integer>create())
                // 对聚合后的结果进行处理
                .apply(MapElements.into(TypeDescriptors.kvs(TypeDescriptors.strings(), TypeDescriptors.integers())).via(e -> {
                    Iterable<Integer> value = e.getValue();
                    if (value == null) {
                        throw new NullPointerException();
                    }
                    Iterator<Integer> iterator = e.getValue().iterator();
                    Integer i = 0;
                    while (iterator.hasNext()) {
                        i += iterator.next();
                    }
                    return KV.of(e.getKey(), i*10);
                }))
                // 自定义算子打印结果集
                .apply(ParDo.of(new DoFn<KV<String, Integer>, KV<String, Integer>>() {
                    @ProcessElement
                    public void processElement(ProcessContext context) {
                        // 从管道中取出的每个元素
                        KV<String, Integer> element = context.element();
                        System.out.println(element);
                        if (element != null) {
                            context.output(element);
                        }
                    }
                }));

        // 将结果集写入数据库
        resultCollection.apply(JdbcIO.<KV<String,Integer>>write()
                .withDataSourceConfiguration(JdbcIO.DataSourceConfiguration.create(
                        "com.mysql.jdbc.Driver",
                        options.getConnUrl())
                        .withUsername(options.getUserName())
                        .withPassword(options.getPassword()))
                .withStatement("insert into TestBeamCount values(?,?)")
                .withPreparedStatementSetter(new JdbcIO.PreparedStatementSetter<KV<String, Integer>>() {
                    @Override
                    public void setParameters(KV<String, Integer> element, PreparedStatement preparedStatement) throws Exception {
                        preparedStatement.setString(1,element.getKey());
                        preparedStatement.setInt(2,element.getValue());
                    }
                }));

                p.run().waitUntilFinish();
    }

    public static void main(String[] args) {
        RankCountByClass2.RankCountByClassOptions options =
                PipelineOptionsFactory.fromArgs(args)
                        .withValidation()
                        .as(RankCountByClass2.RankCountByClassOptions.class);

        runCount(options);
    }
}
