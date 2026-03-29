import { useState, useEffect } from 'react';
import { Line, Pie, Bar } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  ArcElement,
  BarElement,
  ChartOptions
} from 'chart.js';

// 注册Chart.js组件
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  ArcElement,
  BarElement
);

// 接口数据TS类型
interface SalesData {
  months: string[];
  series: Array<{ name: string; data: number[] }>;
}

interface ChannelItem {
  name: string;
  销售额: number;
  订单量: number;
  客单价: number;
}

function App() {
  const [salesData, setSalesData] = useState<SalesData>({ months: [], series: [] });
  const [loading, setLoading] = useState(true);
  const [sortKey, setSortKey] = useState("销售额");
  const [sortType, setSortType] = useState<"asc" | "desc">("desc");

  // 接口请求
  const fetchSalesData = async () => {
    setLoading(true);
    const url = 'https://m1.apifoxmock.com/m1/5076419-0-default/api/sales/monthlySales?apifoxToken=kS5RF-4neMuhg_Qvabu40';
    const res = await fetch(url);
    const result = await res.json();
    setSalesData(result.data);
    setLoading(false);
  };

  useEffect(() => {
    fetchSalesData();
    // 本地排序状态
    const key = localStorage.getItem("sortKey");
    const type = localStorage.getItem("sortType");
    if (key) setSortKey(key);
    if (type) setSortType(type as "asc" | "desc");
  }, []);

  // 表格排序
  const channelData: ChannelItem[] = [
    { name: "WPS官网", 销售额: 452000, 订单量: 3200, 客单价: 141 },
    { name: "App Store", 销售额: 385000, 订单量: 2800, 客单价: 137 },
    { name: "天猫旗舰店", 销售额: 324000, 订单量: 2500, 客单价: 129 },
    { name: "京东自营", 销售额: 298000, 订单量: 2100, 客单价: 142 },
    { name: "企业集采", 销售额: 580000, 订单量: 50, 客单价: 11600 },
  ];

  const sortedData = [...channelData].sort((a, b) => {
    const valA = a[sortKey as keyof ChannelItem] as number;
    const valB = b[sortKey as keyof ChannelItem] as number;
    return sortType === "asc" ? valA - valB : valB - valA;
  });

  const handleSort = (key: string) => {
    const newType = sortKey === key ? (sortType === "desc" ? "asc" : "desc") : "desc";
    setSortKey(key);
    setSortType(newType);
    localStorage.setItem("sortKey", key);
    localStorage.setItem("sortType", newType);
  };

  // 图表
  const lineOptions: any = {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
      x: { grid: { color: 'rgba(255,255,255,0.1)' }, ticks: { color: '#fff' } },
      y: { beginAtZero: true, grid: { color: 'rgba(255,255,255,0.1)' }, ticks: { color: '#fff' } }
    },
    plugins: {
      legend: { labels: { color: '#fff' } }
    }
  };

  // 加载
  if (loading) {
    return <div className="container mx-auto p-4 text-center mt-10 text-white">加载中...</div>;
  }

  // 页面渲染
  return (
    <div className="container mx-auto p-4 bg-[#0F172A] min-h-screen">
      <h1 className="text-2xl font-bold mb-6 text-center text-white">WPS会员年度销售大数据</h1>
      
      {/* 网格布局 */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* 折线图*/}
        <div className="bg-[#1E293B] p-4 rounded-lg h-[400px]">
          <h2 className="text-lg font-semibold mb-3 text-white">年度销售额月度趋势</h2>
          <div className="h-[calc(100%-40px)]">
            <Line
              data={{
                labels: salesData.months,
                datasets: salesData.series.map((item, i) => ({
                  label: item.name,
                  data: item.data,
                  borderColor: ["#06B6D4", "#3B82F6", "#8B5CF6", "#F3F4F6"][i],
                  backgroundColor: "transparent"
                }))
              }}
              options={lineOptions}
            />
          </div>
        </div>

        {/* 饼图 */}
        <div className="bg-[#1E293B] p-4 rounded-lg h-[400px]">
          <h2 className="text-lg font-semibold mb-3 text-white">会员年龄分布</h2>
          <div className="h-[calc(100%-40px)]">
            <Pie
              data={{
                labels: ["18-24岁(学生)", "25-30岁(职场新人)", "31-40岁(资深白领)", "40岁以上"],
                datasets: [{
                  data: [77.78, 15, 5, 2.22],
                  backgroundColor: ["#06B6D4", "#3B82F6", "#8B5CF6", "#EAB308"],
                  borderWidth: 2,
                  borderColor: "#1E293B",
                  hoverOffset: 10
                }]
              }}
              options={{
                responsive: true,
                maintainAspectRatio: false,
                cutout: "60%",
                plugins: {
                  tooltip: {
                    callbacks: {
                      label: (context) => `${context.label}: ${context.raw}%`
                    }
                  },
                  legend: {
                    position: "right",
                    labels: { color: '#fff', padding: 20, font: { size: 12 } }
                  }
                }
              }}
            />
          </div>
        </div>

        {/* 柱状图 */}
        <div className="bg-[#1E293B] p-4 rounded-lg h-[400px]">
          <h2 className="text-lg font-semibold mb-3 text-white">会员地域分布</h2>
          {/* 独立图例 */}
          <div className="flex flex-wrap gap-4 mb-3 text-white">
            <div className="flex items-center"><div className="w-4 h-4 bg-[#EF4444] mr-2"></div><span>广东</span></div>
            <div className="flex items-center"><div className="w-4 h-4 bg-[#EAB308] mr-2"></div><span>其他</span></div>
            <div className="flex items-center"><div className="w-4 h-4 bg-[#84CC16] mr-2"></div><span>北京</span></div>
            <div className="flex items-center"><div className="w-4 h-4 bg-[#10B981] mr-2"></div><span>上海</span></div>
            <div className="flex items-center"><div className="w-4 h-4 bg-[#3B82F6] mr-2"></div><span>浙江</span></div>
            <div className="flex items-center"><div className="w-4 h-4 bg-[#8B5CF6] mr-2"></div><span>江苏</span></div>
            <div className="flex items-center"><div className="w-4 h-4 bg-[#EC4899] mr-2"></div><span>山东</span></div>
            <div className="flex items-center"><div className="w-4 h-4 bg-[#F97316] mr-2"></div><span>四川</span></div>
          </div>
          <div className="h-[calc(100%-80px)]">
            <Bar
              data={{
                labels: ["广东","其他","北京","上海","浙江","江苏","山东","四川","湖北","福建"],
                datasets: [{
                  data: [20500,8000,7800,7200,6500,5800,4500,4200,3800,3500],
                  backgroundColor: ["#EF4444","#EAB308","#84CC16","#10B981","#3B82F6","#8B5CF6","#EC4899","#F97316","#9CA3AF","#4B5563"],
                  borderWidth: 0
                }]
              }}
              options={{
                responsive: true,
                maintainAspectRatio: false,
                indexAxis: "y",
                scales: {
                  x: { grid: { color: 'rgba(255,255,255,0.1)' }, ticks: { color: '#fff' } },
                  y: { grid: { color: 'rgba(255,255,255,0.1)' }, ticks: { color: '#fff' } }
                },
                plugins: { legend: { display: false } }
              }}
            />
          </div>
        </div>

        {/* 表格 */}
        <div className="bg-[#1E293B] p-4 rounded-lg h-[400px]">
          <h2 className="text-lg font-semibold mb-3 text-white">销售渠道</h2>
          <div className="h-[calc(100%-40px)] overflow-auto">
            <table className="w-full text-white">
              <thead>
                <tr className="border-b border-gray-700">
                  <th className="p-2">渠道</th>
                  <th className="p-2 cursor-pointer hover:text-blue-400" onClick={() => handleSort("销售额")}>
                    销售额 {sortKey === "销售额" ? (sortType === "desc" ? "↓" : "↑") : ""}
                  </th>
                  <th className="p-2 cursor-pointer hover:text-blue-400" onClick={() => handleSort("订单量")}>
                    订单量 {sortKey === "订单量" ? (sortType === "desc" ? "↓" : "↑") : ""}
                  </th>
                  <th className="p-2 cursor-pointer hover:text-blue-400" onClick={() => handleSort("客单价")}>
                    客单价 {sortKey === "客单价" ? (sortType === "desc" ? "↓" : "↑") : ""}
                  </th>
                </tr>
              </thead>
              <tbody>
                {sortedData.map((item, i) => (
                  <tr key={i} className="border-t border-gray-700 hover:bg-gray-700/20">
                    <td className="p-2">{item.name}</td>
                    <td className="p-2">{item.销售额}</td>
                    <td className="p-2">{item.订单量}</td>
                    <td className="p-2">{item.客单价}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;