# Использование корутин для чтения файла

Тестовая программа создает файл наполненный псевдослучайным контентом (a-z и пробелы). Далее считает количество вхождений букв разными способами и выводит отчет на экран.

## Сравнение времени подсчета для 100МБ файла (в сек):
    * С корутинами - 1.08 сек
    * Без корутин - 4.63 сек
    * PHP - 3 сек
