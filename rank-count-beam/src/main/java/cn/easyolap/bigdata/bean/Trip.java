package cn.easyolap.bigdata.bean;


import org.apache.beam.repackaged.core.org.apache.commons.lang3.builder.ReflectionToStringBuilder;
import org.apache.beam.repackaged.core.org.apache.commons.lang3.builder.ToStringStyle;

import java.io.Serializable;
import java.math.BigDecimal;
import java.time.LocalDate;

public class Trip implements Serializable,Comparable {
    private String id;
    //person id
    private String sId;
    // class id
    private String cId;
    // cc id
    private String ccId;
    // run distance
    private BigDecimal distance;
    // run time
    private Long runTimes;
    // rank
    private int rank;

    // year
    private String year;
    // month
    private String month;
    // week
    private String week;
    // day
    private String day;

    // create Date
    private LocalDate createDate;



    public String getsId() {
        return sId;
    }

    public void setsId(String sId) {
        this.sId = sId;
    }

    public String getcId() {
        return cId;
    }

    public void setcId(String cId) {
        this.cId = cId;
    }

    public String getCcId() {
        return ccId;
    }

    public void setCcId(String ccId) {
        this.ccId = ccId;
    }

    public BigDecimal getDistance() {
        return distance;
    }

    public void setDistance(BigDecimal distance) {
        this.distance = distance;
    }

    public Long getRunTimes() {
        return runTimes;
    }

    public void setRunTimes(Long runTimes) {
        this.runTimes = runTimes;
    }

    public int getRank() {
        return rank;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public void setRank(int rank) {
        this.rank = rank;
    }

    public String getYear() {
        return year;
    }

    public void setYear(String year) {
        this.year = year;
    }

    public String getMonth() {
        return month;
    }

    public void setMonth(String month) {
        this.month = month;
    }

    public String getWeek() {
        return week;
    }

    public void setWeek(String week) {
        this.week = week;
    }

    public String getDay() {
        return day;
    }

    public void setDay(String day) {
        this.day = day;
    }

    public LocalDate getCreateDate() {
        return createDate;
    }

    public void setCreateDate(LocalDate createDate) {
        this.createDate = createDate;
    }

    @Override
    public int compareTo(Object obj) {
        if(!(obj instanceof Trip))
            throw new RuntimeException("不是学生对象");
        Trip s = (Trip)obj;
        return this.distance.compareTo(s.getDistance());
    }

    @Override
    public String toString()  {
        return ReflectionToStringBuilder.toString(this, ToStringStyle.DEFAULT_STYLE);
    }
}
