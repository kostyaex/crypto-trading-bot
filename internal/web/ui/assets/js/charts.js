
function getRandomColor() {
    // Генерируем случайное число от 0 до 0xFFFFFF, переводим в строку в 16-ричном формате,
    // и дополняем слева нулями до 6 символов
    return `#${Math.floor(Math.random() * 0xFFFFFF).toString(16).padStart(6, '0')}`;
}

// получить данные биржи за период
async function fetchMarketDataForPeriod(start, end) {
    // дата вида 2025-07-23T17:59
    // передать нужно как 2025-07-23T17:59:00+08:00
    // конечная дата должна перевестись в конец минуты
    // знак + экранируем на %2B
    formatedStart = start + ":00%2B08:00"
    formatedEnd = end + ":59%2B08:00"

    console.log(formatedStart)
    console.log(formatedEnd)

    try {
        // здесь дата окончания специально не включительно
        const response = await fetch(`/api/resources/market_data?timestamp=gte.${formatedStart}&timestamp=lt.${formatedEnd}`);
        const items = await response.json();
        //console.log(items)
        return items
    } catch (error) {
        console.error('Ошибка при загрузке данных для редактирования:', error);
    }
}

// заполнить график данными биржи за период
async function fillChart_MarketDataForPeriod(start, end) {

    // получаем данные
    const items = await fetchMarketDataForPeriod(start, end)
    
     // Получаем контейнер по ID
     const chartContainer = document.getElementById('chartContainer');
    
     const chartOptions = {
         layout: { textColor: 'black', background: { type: 'solid', color: 'white' } },
         width: 600,
         height: 400,
     };
     const chart = LightweightCharts.createChart(chartContainer, chartOptions);
     chart.applyOptions({
         timeScale: { timeVisible: true }
     });
     
     const candlestickSeries = chart.addSeries(LightweightCharts.CandlestickSeries, {
         upColor: '#26a69a', downColor: '#ef5350', borderVisible: false,
         wickUpColor: '#26a69a', wickDownColor: '#ef5350',
     });
     
     // преобразуем данные
     //console.log(items)
     const convertedArray = items.map(item => ({
        // преобразуем строку вида 2025-07-23T18:13:00+08:00
        time: Math.floor(new Date(item.timestamp) / 1000),
         open: item.open_price,
         high: item.hight_price,
         low: item.low_price,
         close: item.close_price
     }));
     
     candlestickSeries.setData(convertedArray);
     
     chart.timeScale().fitContent();
}

function fillChartForSeries(series) {
    // Получаем контейнер по ID
    const chartContainer = document.getElementById('chartContainer');
    
    const chartOptions = {
        layout: { textColor: 'black', background: { type: 'solid', color: 'white' } },
        width: 600,
        height: 400,
    };
    const chart = LightweightCharts.createChart(chartContainer, chartOptions);
    chart.applyOptions({
        timeScale: { timeVisible: true }
    });
    
    
    //for (let i = 0; i < series.length; i++) {
    //const color = getRandomColor();
    //const areaSeries = chart.addSeries(LightweightCharts.LineSeries, { color: color });
    
    //originalArray = series[i].points
    //const convertedArray = originalArray.map(item => ({
        //// преобразуем time в секунды и копируем остальные поля
    //value: item.value,
    //time: Math.floor(new Date(item.time).getTime() / 1000)
    //}));
    ////     console.log(convertedArray);
    //areaSeries.setData(convertedArray);
    //}
    
    //const color = getRandomColor();
    //const areaSeries = chart.addSeries(LightweightCharts.LineSeries, { color: color });
    // areaSeries.setData([
    //     { value: 1, time: 1740787800 },
    //     { value: 2, time: 1740787920 }
    // ]);
    //const testdata = [{ value: 0, time: 1642425322 }, { value: 8, time: 1642511722 }, { value: 10, time: 1642598122 }, { value: 20, time: 1642684522 }, { value: 3, time: 1642770922 }, { value: 43, time: 1642857322 }, { value: 41, time: 1642943722 }, { value: 43, time: 1643030122 }, { value: 56, time: 1643116522 }, { value: 46, time: 1643202922 }];
    //areaSeries.setData(testdata);
    
    
    const candlestickSeries = chart.addSeries(LightweightCharts.CandlestickSeries, {
        upColor: '#26a69a', downColor: '#ef5350', borderVisible: false,
        wickUpColor: '#26a69a', wickDownColor: '#ef5350',
    });
    
    originalArray = series.marketdata
    const convertedArray = originalArray.map(item => (        {
        time: Math.floor(new Date(item.Timestamp) / 1000),
        open: item.OpenPrice,
        high: Math.max(item.OpenPrice,item.ClosePrice),
        low: Math.min(item.OpenPrice,item.ClosePrice),
        close: item.ClosePrice
    }));
    // {
    //     // преобразуем time в секунды и копируем остальные поля
    //     time: Math.floor(new Date(item.Timestamp) / 1000),
    //     open: item.OpenPrice,
    //     hight: Math.max(item.OpenPrice,item.ClosePrice),
    //     low: Math.min(item.OpenPrice,item.ClosePrice),
    //     close: item.ClosePrice,
    // }

    //console.log(convertedArray)
    
    candlestickSeries.setData(convertedArray);
    
    
    //candlestickSeries.setData([
    //{ time: '2018-12-22', open: 75.16, high: 82.84, low: 36.16, close: 45.72 },
    //{ time: '2018-12-23', open: 45.12, high: 53.90, low: 45.12, close: 48.09 },
    //{ time: '2018-12-24', open: 60.71, high: 60.71, low: 53.39, close: 59.29 },
    //{ time: '2018-12-25', open: 68.26, high: 68.26, low: 59.04, close: 60.50 },
    //{ time: '2018-12-26', open: 67.71, high: 105.85, low: 66.67, close: 91.04 },
    //{ time: '2018-12-27', open: 91.04, high: 121.40, low: 82.70, close: 111.40 },
    //{ time: '2018-12-28', open: 111.51, high: 142.83, low: 103.34, close: 131.25 },
    //{ time: '2018-12-29', open: 131.33, high: 151.17, low: 77.68, close: 96.43 },
    //{ time: '2018-12-30', open: 106.33, high: 110.20, low: 90.39, close: 98.10 },
    //{ time: '2018-12-31', open: 109.87, high: 114.69, low: 85.66, close: 111.26 },
    //]);
    
    chart.timeScale().fitContent();
}


