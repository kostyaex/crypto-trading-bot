
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

function NewChart() {
    return {
        chart : null,
        containerId: "",
        InitChart (containerId) {
            this.containerId = containerId;
            // Получаем контейнер по ID
            //const chartContainer = document.getElementById(containerId);
                
            //this.chart = Plotly.newPlot(containerId);
        },

        updateForBacktesting(backtesting) {

            const traces = []; // подготовленные данные для графика

            // Серии
            backtesting.series_list.forEach(series => {
                if (series.points.length < 5) {
                    return
                }

                const timeAxis = series.points.map(item => (new Date(item.time)));
                const priceAxis = series.points.map(item => (item.value));

                traces.push({ x: timeAxis, y: priceAxis, type: 'scatter', color: getRandomColor() })
            
            });

            
            var layout = { title: 'Серии' };
            Plotly.react(this.containerId, traces, layout);
            
    
            // Торговые свечи
            // const candlestickSeries = this.chart.addSeries(LightweightCharts.CandlestickSeries, {
            //     upColor: '#26a69a', downColor: '#ef5350', borderVisible: false,
            //     wickUpColor: '#26a69a', wickDownColor: '#ef5350',
            // });
            // this.seriesList.push(candlestickSeries);
            
			// const delta = 1
            // originalArray = backtesting.clustered_marketdata
            // const convertedArray = originalArray.map(item => (        {
            //     time: Math.floor(new Date(item.Timestamp) / 1000),
            //     open: item.ClusterPrice - delta,
            //     high: item.ClusterPrice + delta,
            //     low: item.ClusterPrice - delta,
            //     close: item.ClusterPrice + delta
            // }));

            //candlestickSeries.setData(convertedArray);
        
            //console.log(convertedArray)
              
            //this.chart.timeScale().fitContent();
        }
    }
}





