SELECT * FROM student;
SELECT * FROM teacher;
SELECT * FROM course;
SELECT * FROM sc;

-- https://hackmd.io/@_7vFEnkKTve5g-aFhT8EvQ/Sy-H0QeWr

--1
select * from student s limit 10;

--2
select min(score), avg(score), sum(score)
from sc s;

--3 想詢問 join 和 sub sql 哪個好
select count(c.cno)
from course c 
where c.tno = (select t.tno  from teacher t where t.tname='諶燕');

--通常 count 會用 * 符號
select count(*)
from course c 
join teacher t on c.tno = t.tno and t.tname = '諶燕';

--4
select t.tname, count(c.cno)
from course c 
join teacher t on c.tno = t.tno
group by t.tno 

--5
select *
from student s 
where s.sname like '張%';

--6
select s.sno, sc.score 
from student s 
join sc on s.sno = sc.sno 
join course c on sc.cno = c.cno 
where c.cname = 'Oracle' and sc.score < 60;

--7
select s.sno, s.sname, c.cname 
from course c 
join sc on c.cno = sc.cno 
join student s on s.sno = sc.sno
order by s.sno;

--8 想問 join 表的順序 有什麼條件嗎? 資料量小的放外層?
select s.sno, c.cname, sc.score 
from course c 
join sc on c.cno = sc.cno and sc.score >= 70 
join student s on sc.sno = s.sno
order by s.sno, sc.score desc; 

--9
select s.sno, c.cno, c.cname, sc.score 
from course c 
join sc on c.cno = sc.cno and sc.score < 60 
join student s on sc.sno = s.sno
order by c.cno desc, sc.score;

--10 查詢沒學過”諶燕”老師講授的任一門課程的學號,學生姓名
--注意, 要重複看這題
--預期作法
--先求特定老師教授的課程
--在看分數表 誰沒上過這些課
--但發現作不出來
--
--換另一個想法
--上過特定老師課程的學生有哪些
--在從學生表查詢 誰不再上面

-- 錯誤作法1
--select *
--from student s join sc on s.sno = sc.sno join course c on sc.cno = c.cno left join teacher t on c.tno = t.tno

-- 錯誤作法2
--select t.tname , c.cname , sc.sno, s.sno 
--from teacher t join course c on t.tno = c.tno join sc on c.cno = sc.cno
--right join student s on sc.sno = s.sno;

-- 錯誤作法3
--這個是求 排除特定老師的課程
--學生上了什麼課
--select sc.cno ,sc.sno ,s.sname 
--from sc join student s on sc.sno = s.sno 
--where sc.cno not in (
--	select c.cno 
--	from course c join teacher t on c.tno = t.tno and t.tname = '諶燕'
--) 

-- 錯誤作法4
--select *
--from course c join teacher t on c.tno = t.tno and t.tname = '諶燕'
--right join (
--	select sc.cno ,sc.sno ,s.sname 
--	from sc join student s on sc.sno = s.sno 
--) st_course on c.cno = st_course.cno

--select *
select s.sno ,s.sname 
from course c 
join teacher t on c.tno = t.tno and t.tname = '諶燕' 
join sc on c.cno = sc.cno 
right join student s on sc.sno = s.sno 
where t.tno is null;

--解答版本
SELECT sno,sname
FROM student
WHERE sno NOT in (
	SELECT DISTINCT sno
	FROM sc 
	LEFT JOIN student USING(sno) 
	LEFT JOIN course USING(cno) 
	LEFT JOIN teacher USING(tno)
	WHERE tname IN ('諶燕')
)

--11
-- ERROR: column "sc.sno" must appear in the GROUP BY clause or be used in an aggregate function
--select sc.sno, avg(sc.score) 
--from sc join (
--	select sc.sno, count(sc.sno) as sc_count
--	from sc 
--	where sc.score < 60
--	group by sc.sno 
--) bad_score on sc.sno = bad_score.sno and bad_score.sc_count >=2

--方法1
select sc.sno, avg(sc.score) 
from sc 
join (
	select sc.sno, count(sc.sno) as sc_count
	from sc 
	where sc.score < 60
	group by sc.sno 
) bad_score on sc.sno = bad_score.sno and bad_score.sc_count >=2
group by sc.sno;

--方法2
select sc.sno, avg(sc.score) 
from sc 
join (
	select sc.sno
	from sc 
	where sc.score < 60
	group by sc.sno
	having count(*) >=2
) bad_score on sc.sno = bad_score.sno
group by sc.sno;

--方法3
select sct.sno, student.sname, sct1.avg_score
From ( 
	SELECT sc.sno 
	FROM sc 
	WHERE sc.score < 60 
	GROUP BY sc.sno 
	HAVING COUNT(*)>= 2 
) sct
JOIN ( 
	SELECT sc.sno, AVG(sc.score) avg_score 
	FROM sc 
	GROUP BY sc.sno
) sct1 ON sct1.sno=sct.sno
JOIN student ON sct1.sno = student.sno;

-- 錯誤作法
-- 這只會把 60分以下的成績 進行平均
--select sno , avg(sc.score)
--from sc
--where sc.score < 60
--group by sc.sno 
--having count(sc.sno) >=2

--12
select *
from sc 
where sc.cno = 'c004' and sc.score < 60
order by sc.score desc

--13 要多次複習 
--要比較同一個表 同一個欄位的大小 如何比較
--可看 20 題

--錯誤錯法 少一個必要條件
--select sc1.sno 
--from sc sc1 join sc sc2 on sc1.cno = 'c001' and sc2.cno = 'c002'
--where sc1.score > sc2.score 

select sc1.sno 
from sc sc1 
join sc sc2 on sc1.sno = sc2.sno and sc1.cno = 'c001' and sc2.cno = 'c002' 
where sc1.score > sc2.score 

--解答提供
SELECT a.sno
FROM sc a,sc b
WHERE a.sno = b.sno AND a.cno='c001' AND b.cno='c002' AND a.score > b.score

--14 查詢平均成績大於60 分的同學的學號和平均成績
--必要複習 子查詢 可以放在哪些位置 
--可以參考 從零開始書籍 5-27

--錯誤作法
--以下是列出
--高於自己平均分數有哪些科目
--select *
--from sc sc1 join student s on sc1.sno = s.sno 
--where sc1.score > (
--	select avg(sc2.score) 
--	from sc sc2 
--	where sc1.sno = sc2.sno 
--	group by sc2.sno 
--)
--order by sc1.sno 

select sc.sno , avg(sc.score) 
from sc
group by sc.sno 
having avg(sc.score) > 60

--15

select sc.sno , count(sc.cno) as course_count, sum(sc.score) 
from student s 
join sc on s.sno = sc.sno 
group by sc.sno 

--16
select count(t.tno)
from teacher t 
where t.tname like '劉%'

--17 查詢只學”諶燕”老師所教的課的同學的學號:姓名
--注意 DISTINCT 的用法
--錯誤作法 可以和 19題進行比較
--select distinct s.sno ,s.sname 
--from student s join sc on s.sno = sc.sno join course c on sc.cno = c.cno  join teacher t on c.tno = t.tno and t.tname = '諶燕'

select *
from sc , count(*)
join course c using(sc.cno)
join ( 
	select sc.sno, c.tno, count(*)
	from sc
	join course c using(cno)
	group by (sc.sno, c.tno)

) 

select *
from student s 
join sc using(sno)
join course c using(cno)
join teacher t using(tno)
where s.sno in (
	--利用 count 求出學生選課幾位老師
	--最後得到只選一位老師進行修課的學生名單
	select st_te.sno
	from (
		--先求 所有學生選了每個老師幾門課程
		--如果選了一個學生選了兩個老師的課
		--那就會有兩條紀錄
		select sc.sno, c.tno, count(*)
		from sc
		join course c using(cno)
		join teacher t using(tno)
		group by sc.sno, c.tno
	) st_te
	group by st_te.sno
	having count(*) = 1
)
--and t.tname = '諶燕'


--以下三個sql
--參考的維度是相同的
--可以比較看看
--
--也應該想想 連續 join 的時候
--是為了什麼原因 才選擇用某個 table 當作第一個 table
--
--sql1
--select distinct s.sno, tc.tname
--from sc s 
--left join (
--    select t.tname, c.cno 
--    from course c 
--    left join teacher t on c.tno = t.tno
--) tc on s.cno = tc.cno;
--
--sql2
--select distinct s.sno, t.tname
--from sc s 
--join course c using(cno)
--join teacher t using(tno)
--
--sql3
--select sc.sno, c.tno, count(*)
--from sc
--join course c using(cno)
--join teacher t using(tno)
--group by sc.sno, c.tno

--錯誤答案 
--但值得思考 group by 用 一個 跟 兩個 欄位的區別
--group by sc.sno vs group by sc.sno, c.tno
--
--select sc.sno, c.tno, count(*)
--from sc
--join course c using(cno)
--group by sc.sno, c.tno

--18 查詢學過”c001″並且也學過編號”c002″課程的同學的學號.姓名

--錯誤作法1 仔細想想條件
--不可能同時有一個表的紀錄 同時滿足 等於 c001 和 c002 
--select *
--from student s 
--join sc on s.sno = sc.sno and sc.cno = 'c001' and sc.cno = 'c002'

--錯誤作法2 仔細想想條件
--select *
--from student s 
--join sc on s.sno = sc.sno and sc.cno = 'c001' or sc.cno = 'c002'

--錯誤作法 這只能找出 學過 'c001' or 'c002'
--注意 in 的用法
--用小括號表示 list
--select *
--from student s join sc on s.sno = sc.sno and (sc.cno in  ('c001', 'c002'))

--解答
SELECT sno, sname
FROM student
WHERE sno in (
	SELECT a.sno
	FROM sc a,sc b 
	WHERE a.sno=b.sno AND a.cno='c001' AND b.cno='c002'
)

--19 查詢學過”諶燕”老師所教的所有課的同學的學號:姓名
--可和 17 比較語意上的差別
select distinct s.sno ,s.sname 
from student s 
join sc on s.sno = sc.sno 
join course c on sc.cno = c.cno  
join teacher t on c.tno = t.tno 
where t.tname = '諶燕'

SELECT sno, sname
FROM sc
LEFT JOIN student USING(sno) 
LEFT JOIN course USING(cno) 
LEFT JOIN teacher USING(tno)
WHERE cno in (
	SELECT cno 
	FROM course 
	LEFT JOIN teacher USING(tno) 
	WHERE tname='諶燕'
	)
GROUP BY sno, sname
HAVING COUNT(*) = (SELECT COUNT(*) FROM course LEFT JOIN teacher USING(tno) WHERE tname='諶燕')

--20 查詢課程編號”c004″的成績比課程編號”c001″和”c002″課程低的所有同學的學號.姓名
--可看 13 題

--ERROR: more than one row returned by a subquery used as an expression
--想想 為什麼這個 subquery 無法執行
--和 14 的差異?
--select *
--from student s 
--join sc sc1 on s.sno = sc1.sno and sc1.cno = 'c004'
--where sc1.score < (
--	select sc2.score 
--	from sc sc2
--	where sc1.sno = sc2.sno and (sc2.cno in ('c001', 'c002'))
--)

--錯誤作法 無法比較
-- (sc2.cno = 'c001' or sc2.cno = 'c002') 一個 join 裡面用 or
-- 只能增加紀錄, 無法擴展欄位
-- 要比較同一個表 同一個欄位的大小 需要擴展欄位!
--
--select *
--from sc sc1 
--join sc sc2 on sc1.sno = sc2.sno and sc1.cno = 'c004' and (sc2.cno = 'c001' or sc2.cno = 'c002')
--join student s on sc1.sno = s.sno 

--思考一個問題, 如果要比較的課程更多, 是否要 join 更多次?
select *
from sc sc1 
join sc sc2 on sc1.sno = sc2.sno 
join sc sc3 on sc1.sno = sc3.sno 
join student s on sc1.sno = s.sno 
where sc1.cno = 'c004' and sc2.cno = 'c001' and sc3.cno = 'c002' and sc1.score < sc2.score and sc1.score < sc3.score

select *
from sc sc1 
join sc sc2 using(sno)
join sc sc3 using(sno)
join student s using(sno)
where sc1.cno = 'c004' and sc2.cno = 'c001' and sc3.cno = 'c002' and sc1.score < sc2.score and sc1.score < sc3.score

--21 查詢所修課程其成績都小於60分的同學的學號.姓名

--錯誤
--這是查詢不及格課程的學生名單
--select s.sname , sc.score 
--from sc
--join student s using(sno)
--where sc.score < 60

--解答
--關聯子查詢考察: 14, 21, 29
--
--三個關鍵字可以修改比較運算符
--ALL、ANY和SOME
--作用於比較運算符和子查詢之間，作用類似EXISTS
--ALL：是所有，表示全部都滿足才返回true；
--ANY/SOME：是任意一個 ，表示有任何一個滿足就返回true。
SELECT DISTINCT sno
FROM sc x
WHERE 60 > ALL(
	SELECT score 
	FROM sc y 
	WHERE x.sno=y.sno
)all

--22
select *
from student s 
where s.sno not in (
	select distinct sc.sno 
	from sc
)

--23
select distinct sc1.sno , s.sname 
from sc sc1
join sc sc2 using(cno)
join student s on s.sno = sc1.sno 
where sc2.sno = 's001' and sc1.sno != 's001'

-- 解答
SELECT DISTINCT sc.sno,student.sname
FROM sc 
RIGHT JOIN student USING(sno)
WHERE cno IN(
	SELECT DISTINCT cno
	FROM sc
	WHERE sno='s001'
)
AND sc.sno <> 's001'

--24 查詢跟學號為”s005″所修課程完全一樣的同學的學號和姓名
--可以思考 如果沒先出現 23題
--自己會怎麼寫

--錯誤作法
--再次想想 跟自己的欄位相比 進行 join
--會回傳什麼值
--select *
--from (
--	select sc1.sno, count(*) as co_count
--	from sc sc1
--	join sc sc2 using(cno)
--	group by sc1.sno 
--	) st_co1 
--join (
--	select sc1.sno, count(*) as co_count
--	from sc sc1
--	join sc sc2 using(cno)
--	group by sc1.sno 
--)  st_co2 on st_co1.sno = st_co2.sno
--join student s on s.sno = st_co2.sno
--where st_co1.sno = 's005' and st_co1.co_count = st_co2.co_count

--錯誤答案
--以下 sql 只能說明 上課人數相同
--不代表上同樣的課
--
--雖然錯了, 但可以激發思考
--20題: 想要比較同一筆紀錄中, 不同欄位的相關性, 利用 join
--24題: 想要比較不同紀錄, 相同欄位的相關性, 應該使用 cross join
--
--select st_co2.sno
--from (
--	--找出同一門課
--	--包含自己的學生數量
--	select sc1.sno, count(*) as co_count
--	from sc sc1
--	join sc sc2 using(cno)
--	group by sc1.sno
--	) st_co1 
--cross join (
--	select sc1.sno, count(*) as co_count
--	from sc sc1
--	join sc sc2 using(cno)
--	group by sc1.sno
--)  st_co2
--where st_co1.sno = 's005' 
--and st_co1.co_count = st_co2.co_count
--and st_co2.sno <> 's005'

--解答1
SELECT x.sno, t.sname
FROM sc x 
LEFT JOIN student t USING(sno)
WHERE cno in (SELECT cno FROM sc WHERE sno = 's005') -- 其中一門課 和 s005 一樣就會入選
GROUP BY x.sno, t.sname 
HAVING COUNT(cno) = (SELECT COUNT(*) FROM sc WHERE sno = 's005') -- 課程數量相同
AND (SELECT COUNT(*) FROM sc WHERE sno = 's005') = ALL(SELECT COUNT(cno) FROM sc y WHERE x.sno=y.sno ) -- 確保真的都選一樣的課
AND x.sno <> 's005';


--25
select sc.cno, min(sc.score), max(sc.score) 
from sc
group by sc.cno 

--26 按各科平均成績和及格率的百分數 照平均從低到高顯示

--沒有想法怎麼求 及格率
--算及格率
--若不存在有人集合
--不會回傳0
--而是不回傳任何資訊
--
--select sc.cno, count(*)
--from sc
--where sc.score > 60
--group by sc.cno

--解答
--COALESCE() Function: Return the first non-null value in a list
--CONCAT()函数可将两个或多个字符串连接成一个字符串
SELECT cno, co_avg.avg_score , CONCAT(COALESCE(pcount,0)/tcount*100,'%') AS passing
FROM (
	SELECT cno, avg(score) AS avg_score FROM sc GROUP BY cno
) AS co_avg
LEFT JOIN (
	SELECT cno, COUNT(*) AS pcount FROM sc WHERE score > 60 GROUP BY cno
) AS pass USING(cno)
LEFT JOIN (
	SELECT cno, COUNT(*) AS tcount FROM sc GROUP BY cno
) AS total USING(cno)
ORDER BY co_avg.avg_score

--27
select sc.cno, t.tname, avg(sc.score) as avg_score
from sc
join course c using(cno)
join teacher t using(tno)
group by sc.cno, t.tname  
order by avg_score desc

--28 統計列印各科成績,各分數段人數:課程ID,課程名稱,verygood[100-86], good[85-71], bad[<60]

--解答1
SELECT cno,cname, COALESCE(verygoodc,0) as verygood, COALESCE(goodc,0) as good, COALESCE(badc,0) as bad
FROM sc 
LEFT JOIN (
	SELECT cno,COUNT(*) verygoodc FROM sc WHERE score BETWEEN 86 AND 100 GROUP BY cno
) AS verygoodsc USING(cno)
LEFT JOIN (
	SELECT cno,COUNT(*) goodc FROM sc WHERE score BETWEEN 71 AND 85 GROUP BY cno
) AS goodsc USING(cno)
LEFT JOIN (
	SELECT cno,COUNT(*) badc FROM sc WHERE score < 60 GROUP BY cno
) AS badsc USING(cno)
LEFT JOIN course USING(cno)

--解答2
SELECT sc.cno,course.cname,
SUM(CASE WHEN sc.score BETWEEN 86 AND 100 THEN 1 ELSE 0 END) as verygood,
SUM(CASE WHEN sc.score BETWEEN 71 AND 85 THEN 1 ELSE 0 END) as good,
SUM(CASE WHEN sc.score < 60 THEN 1 ELSE 0 END) as bad
FROM sc,course
WHERE sc.cno=course.cno
GROUP BY cno

--29 查詢各科成績前三名的記錄:(不考慮成績並列情況)
--看完解答還是覺得想不透
--關聯子查詢考察: 14, 21, 29

--對於取前三個 沒有想法
--錯誤作法
--select *
--from sc sc1
--join sc sc2 using(cno)

--解答1
--單純用 join 無法作到 同時比較關係 又 紀錄數量
SELECT *
FROM sc x
WHERE (
	-- 比 x 分數大的學生有幾人
	-- 0 人 表示 x 是 第一名, x 沒輸任何人
	-- 1 人 表示 x 是 第二名, x 輸1人
	SELECT COUNT(*) FROM sc y WHERE x.cno = y.cno AND x.score < y.score
) < 3 
ORDER BY cno,score desc

--解答2
--join 後可能出現 NULL 欄位
--所以 count 一定要指定 計算的欄位, 不可用 count(*)
--不然答案會錯
--
--HAVING子句中能够使用三种要素：
--1. 常数
--2. 聚合函数
--3. GROUP BY子句中指定的列名(聚合建)
SELECT a.*, COUNT(b.score) +1 AS ranking
FROM sc AS a 
LEFT JOIN sc AS b ON a.cno = b.cno AND a.score<b.score
GROUP BY a.cno ,a.sno
HAVING count(b.score)+1 <= 3
ORDER BY a.cno,ranking;

--測試用途
SELECT *
FROM sc AS a 
LEFT JOIN sc AS b ON a.cno = b.cno AND a.score<b.score

--window func
select sno, cno, score
from (
	select sc.sno, sc.cno, sc.score , dense_rank() over( partition by sc.cno order by sc.score desc ) as course_rank
	from sc
) as sc_rank
where course_rank <= 3

--30
select sc.cno , count(*) as st_count
from sc
group by sc.cno 




