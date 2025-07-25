
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
        seriesList: [],
        InitChart (containerId) {
            // Получаем контейнер по ID
            const chartContainer = document.getElementById(containerId);
                
            const chartOptions = {
                layout: { textColor: 'black', background: { type: 'solid', color: 'white' } },
                //width: 600,
                height: 400,
            };
            this.chart = LightweightCharts.createChart(chartContainer, chartOptions);
            this.chart.applyOptions({
                timeScale: { timeVisible: true }
            });
        },

        updateForBacktesting(backtesting) {

            // Удаляем все существующие серии
            this.seriesList.forEach(series => this.chart.removeSeries(series));
    
            const candlestickSeries = this.chart.addSeries(LightweightCharts.CandlestickSeries, {
                upColor: '#26a69a', downColor: '#ef5350', borderVisible: false,
                wickUpColor: '#26a69a', wickDownColor: '#ef5350',
            });
            this.seriesList.push(candlestickSeries);
            
            originalArray = backtesting.clustered_marketdata
            const convertedArray = originalArray.map(item => (        {
                time: Math.floor(new Date(item.Timestamp) / 1000),
                open: item.ClusterPrice-50,
                high: item.ClusterPrice+50,
                low: item.ClusterPrice-50,
                close: item.ClusterPrice+50
            }));
        
            //console.log(convertedArray)
            
            candlestickSeries.setData(convertedArray);
            
            
            this.chart.timeScale().fitContent();
        }
    }
}





