function series() {
    return {
        seriesUrls: [],
        series: [],
        status: "",
        async loadSeriesUrls() {
            try {
                const response = await fetch('/api/seriesdumpslist', {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(this.form),
                });

                this.seriesUrls = await response.json();
                this.status = "Серии загружены"

            } catch (error) {
                this.status = 'Ошибка при отправке данных:'+error
                console.error('Ошибка при отправке данных:', error);
            }
        },
        async loadSeries(url) {
            try {
                const response = await fetch('/data'+url, {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(this.form),
                });

                this.series = await response.json();
                //console.log(data)
                //fillChartForSeries(data)
                this.status = "Загружены данные серии "+url

            } catch (error) {
                this.status = 'Ошибка при отправке данных:'+error
                console.error('Ошибка при отправке данных:', error);
            }
        }
    };
}

