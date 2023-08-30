FROM debian

# переменная для запуска с веб-интерфесом по дефолту
# чтобы запустить без - нужно переопределить на любое другое значение или пустое
ENV WITH_UI=yes

RUN apt-get update
RUN apt-get install 

RUN apt install -y golang 

# устанавливаем воркдир
WORKDIR /build  

COPY . .

WORKDIR /build/src

RUN go build -o app main.go


#CMD ["./app", "${WITH_UI}"]
CMD ./app ${WITH_UI}

