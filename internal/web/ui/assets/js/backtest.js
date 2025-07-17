function backtests() {
    return {
        urls: [],
        backtests: [],
        status: "",
        async backtestsUrls() {
            try {
                const response = await fetch('/api/backtestsdumpslist', {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(this.form),
                });

                this.urls = await response.json();
                this.status = "Результаты загружены"

            } catch (error) {
                this.status = 'Ошибка при отправке данных:'+error
                console.error('Ошибка при отправке данных:', error);
            }
        },
        async loadBacktest(url) {
            try {
                const response = await fetch('/data'+url, {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(this.form),
                });

                this.backtests = await response.json();
                //console.log(data)
                //fillChartForSeries(data)
                this.status = "Загружены данные бектеста "+url

            } catch (error) {
                this.status = 'Ошибка при получении данных:'+error
                console.error('Ошибка при получении данных:', error);
            }
        }
    };
}

