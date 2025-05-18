-- Скрипт для создания тестовых данных в БД FurniSwap
-- Создает 10 пользователей с объявлениями, покупками, избранным и чатами

-- Очистка существующих данных (опционально)
TRUNCATE TABLE messages CASCADE;
TRUNCATE TABLE chats CASCADE;
TRUNCATE TABLE favorites CASCADE;
TRUNCATE TABLE purchases CASCADE;
TRUNCATE TABLE listing_images CASCADE;
TRUNCATE TABLE listings CASCADE;
TRUNCATE TABLE two_factor_codes CASCADE;
TRUNCATE TABLE users CASCADE;

-- Сброс ID последовательностей
ALTER SEQUENCE users_id_seq RESTART WITH 1;
ALTER SEQUENCE listings_id_seq RESTART WITH 1;
ALTER SEQUENCE listing_images_id_seq RESTART WITH 1;
ALTER SEQUENCE chats_id_seq RESTART WITH 1;
ALTER SEQUENCE messages_id_seq RESTART WITH 1;
ALTER SEQUENCE favorites_id_seq RESTART WITH 1;
ALTER SEQUENCE purchases_id_seq RESTART WITH 1;

-- Создание 10 пользователей (хеш пароля 'password123')
INSERT INTO users (email, password_hash, name, last_name, city, avatar, is_verified, created_at) VALUES
('user1@example.com', '$2a$10$ASTWDJDs52jll8XgyZOkS.43W.u.L2BdtV1N/RiwRMhHvwnH5hfz.', 'Анна', 'Иванова', 'Москва', 'https://randomuser.me/api/portraits/women/1.jpg', true, NOW() - INTERVAL '30 days'),
('user2@example.com', '$2a$10$ASTWDJDs52jll8XgyZOkS.43W.u.L2BdtV1N/RiwRMhHvwnH5hfz.', 'Сергей', 'Петров', 'Санкт-Петербург', 'https://randomuser.me/api/portraits/men/1.jpg', true, NOW() - INTERVAL '29 days'),
('user3@example.com', '$2a$10$ASTWDJDs52jll8XgyZOkS.43W.u.L2BdtV1N/RiwRMhHvwnH5hfz.', 'Екатерина', 'Смирнова', 'Новосибирск', 'https://randomuser.me/api/portraits/women/2.jpg', true, NOW() - INTERVAL '28 days'),
('user4@example.com', '$2a$10$ASTWDJDs52jll8XgyZOkS.43W.u.L2BdtV1N/RiwRMhHvwnH5hfz.', 'Дмитрий', 'Волков', 'Екатеринбург', 'https://randomuser.me/api/portraits/men/2.jpg', true, NOW() - INTERVAL '27 days'),
('user5@example.com', '$2a$10$ASTWDJDs52jll8XgyZOkS.43W.u.L2BdtV1N/RiwRMhHvwnH5hfz.', 'Мария', 'Соколова', 'Казань', 'https://randomuser.me/api/portraits/women/3.jpg', true, NOW() - INTERVAL '26 days'),
('user6@example.com', '$2a$10$ASTWDJDs52jll8XgyZOkS.43W.u.L2BdtV1N/RiwRMhHvwnH5hfz.', 'Алексей', 'Морозов', 'Самара', 'https://randomuser.me/api/portraits/men/3.jpg', true, NOW() - INTERVAL '25 days'),
('user7@example.com', '$2a$10$ASTWDJDs52jll8XgyZOkS.43W.u.L2BdtV1N/RiwRMhHvwnH5hfz.', 'Ольга', 'Козлова', 'Ростов-на-Дону', 'https://randomuser.me/api/portraits/women/4.jpg', true, NOW() - INTERVAL '24 days'),
('user8@example.com', '$2a$10$ASTWDJDs52jll8XgyZOkS.43W.u.L2BdtV1N/RiwRMhHvwnH5hfz.', 'Павел', 'Лебедев', 'Уфа', 'https://randomuser.me/api/portraits/men/4.jpg', true, NOW() - INTERVAL '23 days'),
('user9@example.com', '$2a$10$ASTWDJDs52jll8XgyZOkS.43W.u.L2BdtV1N/RiwRMhHvwnH5hfz.', 'Наталья', 'Новикова', 'Краснодар', 'https://randomuser.me/api/portraits/women/5.jpg', true, NOW() - INTERVAL '22 days'),
('user10@example.com', '$2a$10$ASTWDJDs52jll8XgyZOkS.43W.u.L2BdtV1N/RiwRMhHvwnH5hfz.', 'Игорь', 'Федоров', 'Воронеж', 'https://randomuser.me/api/portraits/men/5.jpg', true, NOW() - INTERVAL '21 days');

-- Создание объявлений (по 3-5 для каждого пользователя)
INSERT INTO listings (user_id, title, description, price, condition, city, category_id, status, created_at, updated_at) VALUES
-- Пользователь 1
(1, 'Кожаный диван', 'Кожаный диван в отличном состоянии, использовался меньше года', 25000, 'хорошее', 'Москва', 1, 'active', NOW() - INTERVAL '29 days', NOW() - INTERVAL '29 days'),
(1, 'Обеденный стол', 'Деревянный обеденный стол, вместимость 6 человек', 12000, 'среднее', 'Москва', 2, 'active', NOW() - INTERVAL '28 days', NOW() - INTERVAL '28 days'),
(1, 'Кровать с матрасом', 'Двуспальная кровать с ортопедическим матрасом', 18000, 'хорошее', 'Москва', 4, 'active', NOW() - INTERVAL '27 days', NOW() - INTERVAL '27 days'),
(1, 'Книжный шкаф', 'Вместительный книжный шкаф из натурального дерева', 8500, 'среднее', 'Москва', 3, 'active', NOW() - INTERVAL '26 days', NOW() - INTERVAL '26 days'),

-- Пользователь 2
(2, 'Кресло-качалка', 'Удобное кресло-качалка для отдыха', 7500, 'хорошее', 'Санкт-Петербург', 1, 'active', NOW() - INTERVAL '28 days', NOW() - INTERVAL '28 days'),
(2, 'Компьютерный стол', 'Эргономичный компьютерный стол с полками', 6000, 'среднее', 'Санкт-Петербург', 2, 'active', NOW() - INTERVAL '27 days', NOW() - INTERVAL '27 days'),
(2, 'Комод для одежды', 'Комод с 5 ящиками, цвет - венге', 9000, 'хорошее', 'Санкт-Петербург', 3, 'active', NOW() - INTERVAL '26 days', NOW() - INTERVAL '26 days'),

-- Пользователь 3
(3, 'Диван-кровать', 'Раскладной диван-кровать, механизм "еврокнижка"', 15000, 'среднее', 'Новосибирск', 1, 'active', NOW() - INTERVAL '27 days', NOW() - INTERVAL '27 days'),
(3, 'Журнальный столик', 'Стильный журнальный столик из стекла и металла', 4500, 'хорошее', 'Новосибирск', 2, 'active', NOW() - INTERVAL '26 days', NOW() - INTERVAL '26 days'),
(3, 'Шкаф-купе', 'Вместительный шкаф-купе с зеркальными дверями', 22000, 'хорошее', 'Новосибирск', 3, 'active', NOW() - INTERVAL '25 days', NOW() - INTERVAL '25 days'),
(3, 'Пуфик для прихожей', 'Практичный пуфик с отсеком для хранения', 3000, 'новое', 'Новосибирск', 5, 'active', NOW() - INTERVAL '24 days', NOW() - INTERVAL '24 days'),

-- Пользователь 4
(4, 'Мягкое кресло', 'Уютное мягкое кресло, обивка - микровелюр', 8200, 'хорошее', 'Екатеринбург', 1, 'active', NOW() - INTERVAL '26 days', NOW() - INTERVAL '26 days'),
(4, 'Барный стул', 'Высокий барный стул, регулируемая высота', 3500, 'среднее', 'Екатеринбург', 2, 'active', NOW() - INTERVAL '25 days', NOW() - INTERVAL '25 days'),
(4, 'Детская кроватка', 'Детская кроватка с матрасом, бортиками и ящиком', 7000, 'хорошее', 'Екатеринбург', 4, 'active', NOW() - INTERVAL '24 days', NOW() - INTERVAL '24 days'),

-- Пользователь 5
(5, 'Угловой диван', 'Большой угловой диван, раскладывается', 28000, 'хорошее', 'Казань', 1, 'active', NOW() - INTERVAL '25 days', NOW() - INTERVAL '25 days'),
(5, 'Обеденные стулья', 'Комплект из 4 обеденных стульев, дерево и ткань', 10000, 'среднее', 'Казань', 2, 'active', NOW() - INTERVAL '24 days', NOW() - INTERVAL '24 days'),
(5, 'Тумба под ТВ', 'Современная тумба под ТВ с ящиками', 6500, 'хорошее', 'Казань', 3, 'active', NOW() - INTERVAL '23 days', NOW() - INTERVAL '23 days'),
(5, 'Полка настенная', 'Декоративная полка для книг и сувениров', 2200, 'хорошее', 'Казань', 5, 'active', NOW() - INTERVAL '22 days', NOW() - INTERVAL '22 days'),

-- Пользователь 6
(6, 'Офисное кресло', 'Эргономичное офисное кресло на колесиках', 5500, 'среднее', 'Самара', 1, 'active', NOW() - INTERVAL '24 days', NOW() - INTERVAL '24 days'),
(6, 'Письменный стол', 'Классический письменный стол с ящиками', 7800, 'хорошее', 'Самара', 2, 'active', NOW() - INTERVAL '23 days', NOW() - INTERVAL '23 days'),
(6, 'Шкаф для одежды', 'Двухстворчатый шкаф для одежды, цвет - белый', 11000, 'хорошее', 'Самара', 3, 'active', NOW() - INTERVAL '22 days', NOW() - INTERVAL '22 days'),

-- Пользователь 7
(7, 'Мягкий уголок', 'Мягкий уголок с механизмом трансформации', 20000, 'среднее', 'Ростов-на-Дону', 1, 'active', NOW() - INTERVAL '23 days', NOW() - INTERVAL '23 days'),
(7, 'Стеклянный стол', 'Стильный обеденный стол со стеклянной столешницей', 9500, 'хорошее', 'Ростов-на-Дону', 2, 'active', NOW() - INTERVAL '22 days', NOW() - INTERVAL '22 days'),
(7, 'Комод с зеркалом', 'Комод с большим зеркалом для спальни', 13500, 'хорошее', 'Ростов-на-Дону', 3, 'active', NOW() - INTERVAL '21 days', NOW() - INTERVAL '21 days'),
(7, 'Подставка для цветов', 'Металлическая подставка для комнатных растений', 1800, 'новое', 'Ростов-на-Дону', 5, 'active', NOW() - INTERVAL '20 days', NOW() - INTERVAL '20 days'),

-- Пользователь 8
(8, 'Кожаное кресло', 'Кожаное кресло для гостиной, цвет - коричневый', 12000, 'хорошее', 'Уфа', 1, 'active', NOW() - INTERVAL '22 days', NOW() - INTERVAL '22 days'),
(8, 'Складной стол', 'Компактный складной стол для дачи', 3000, 'среднее', 'Уфа', 2, 'active', NOW() - INTERVAL '21 days', NOW() - INTERVAL '21 days'),
(8, 'Двуспальная кровать', 'Двуспальная кровать с подъемным механизмом', 17000, 'хорошее', 'Уфа', 4, 'active', NOW() - INTERVAL '20 days', NOW() - INTERVAL '20 days'),

-- Пользователь 9
(9, 'Диван для офиса', 'Компактный диван для офиса или приемной', 14000, 'хорошее', 'Краснодар', 1, 'active', NOW() - INTERVAL '21 days', NOW() - INTERVAL '21 days'),
(9, 'Стулья для кухни', 'Набор из 4 стульев для кухни, металл и пластик', 6000, 'среднее', 'Краснодар', 2, 'active', NOW() - INTERVAL '20 days', NOW() - INTERVAL '20 days'),
(9, 'Шкаф для посуды', 'Шкаф-витрина для посуды и декора', 9800, 'хорошее', 'Краснодар', 3, 'active', NOW() - INTERVAL '19 days', NOW() - INTERVAL '19 days'),
(9, 'Вешалка напольная', 'Деревянная напольная вешалка для одежды', 2500, 'хорошее', 'Краснодар', 5, 'active', NOW() - INTERVAL '18 days', NOW() - INTERVAL '18 days'),

-- Пользователь 10
(10, 'Кресло-мешок', 'Удобное кресло-мешок, наполнитель - пенополистирол', 4000, 'хорошее', 'Воронеж', 1, 'active', NOW() - INTERVAL '20 days', NOW() - INTERVAL '20 days'),
(10, 'Компьютерное кресло', 'Компьютерное кресло с подлокотниками', 6500, 'среднее', 'Воронеж', 1, 'active', NOW() - INTERVAL '19 days', NOW() - INTERVAL '19 days'),
(10, 'Стол-трансформер', 'Стол-трансформер, регулируемая высота', 8500, 'хорошее', 'Воронеж', 2, 'active', NOW() - INTERVAL '18 days', NOW() - INTERVAL '18 days'),
(10, 'Детский шкаф', 'Яркий шкаф для детской комнаты', 7200, 'хорошее', 'Воронеж', 3, 'active', NOW() - INTERVAL '17 days', NOW() - INTERVAL '17 days');

-- Добавление изображений к объявлениям
INSERT INTO listing_images (listing_id, image_path, is_main, created_at) VALUES
-- Изображения для объявлений пользователя 1
(1, 'https://www.restmebel.ru/upload/iblock/54d/50w428cmi9v8jkok7zygufe17xktmyj7.jpg', true, NOW() - INTERVAL '29 days'),
(1, 'https://lalume.ru/image/cache/catalog/finnnavian/35597-29/o1cn018yzyl2260ouiv1tc3_%21%212214781157599-1440x1440.webp', false, NOW() - INTERVAL '29 days'),
(2, 'https://www.mebel-1000.ru/upload/iblock/62d/62d83759f2e175a76cc71b30f1ff0795.jpg', true, NOW() - INTERVAL '28 days'),
(3, 'https://davitamebel.ru/upload/resize_cache/product/841_631_1/957c2fc6-59a4-11ec-aa34-00155db6cc01_im-00001444_kashemir-900-304/photo_00_im-00001444_kashemir-900-304_krovat-s-matrasom.jpg?1744204057', true, NOW() - INTERVAL '27 days'),
(4, 'https://www.shkaf-kupe.ru/cdn/uGECqgptc8qC4mbsDh2bWRowK1yPFIbqc2E3lzh4REQ/rs:fit:0:0/q:100/bg:FFFFFF/exar:1/aG9tZS93bXAvX2FwcC9fYWN0aXZlL2FwcC9zaGthZi1rdXBlLXJ1L2Zyb250ZW5kL3dlYi9zdGF0aWMvMWUxM2M5OTIvOTA1NzUyLmpwZw.webp', true, NOW() - INTERVAL '26 days'),

-- Изображения для объявлений пользователя 2
(5, 'https://deamo-rotang.ru/images/product_images/popup_images/1332_10_3z8.jpg', true, NOW() - INTERVAL '28 days'),
(6, 'https://avatars.mds.yandex.net/get-mpic/16105123/2a000001969fce90de0fb6d82be94bf4399e/orig', true, NOW() - INTERVAL '27 days'),
(7, 'https://ir.ozone.ru/s3/multimedia-1-m/c400/7286215990.jpg', true, NOW() - INTERVAL '26 days'),

-- Изображения для объявлений пользователя 3
(8, 'https://www.stolplit.ru/upload/resize_cache/iblock/abc/x9pps84ifks3wisk5xtrhmt4tro2ty09/306_262_0/r0000287378-detail.jpg', true, NOW() - INTERVAL '27 days'),
(9, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQZAFatZ8daAkRiYOlEja-R7VoE5koCFdEj5Q&s', true, NOW() - INTERVAL '26 days'),
(10, 'https://mebelyes.com/wp-content/uploads/2021/05/216-1.jpg', true, NOW() - INTERVAL '25 days'),
(11, 'https://krimmebel.ru/image/catalog/blog/banketka1.jpg', true, NOW() - INTERVAL '24 days'),

-- Изображения для других объявлений
(12, 'https://cidien.lazurit.com/images/unsafe/fit-in/1000x1000/upload.lazurit.com/upload/iblock/647/y60s82q657xrvgeunn38pd2ud2gdbeze/kDZH01-kreslo-Dzhordzh-TITANIUM_01-900.jpg', true, NOW() - INTERVAL '26 days'),
(13, 'https://cdn0.divan.ru/img/v1/fwno3aUTY9nO43D6wbvBR34S9z6dpmxPNxlaHy9cMGE/t:0::0:0/pd:60:60:60:60/rs:fit:1148:720:0:1:ce:0:0/g:ce:0:0/bg:f5f3f1/q:95/czM6Ly9kaXZhbi9wcm9kdWN0LzQzNTM1NzAucG5n.jpg', true, NOW() - INTERVAL '25 days'),
(14, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSFJpGB7ai6_RswEGdUUAJMODr5JkWaopIyHg&s', true, NOW() - INTERVAL '24 days'),
(15, 'https://cdn.domdivanov45.com/files/imgs/ig1111921/uglovoi-divan-lider-sv-1-d-14otp07-npb-590x430.jpg', true, NOW() - INTERVAL '25 days'),
(16, 'https://e-mall.su/wa-data/public/shop/products/26/10/11026/images/34071/34071.970.jpg', true, NOW() - INTERVAL '24 days'),
(17, 'https://cdn.nonton.ru/webp/ad7/ad70bdab7470a7abf546c8f1293916d6/08_101975_Tumba_pod_TV_Skandi_Seryy_SHifer_matovaya_Grin_Grey_Soft_rum2.jpg.webp', true, NOW() - INTERVAL '23 days'),
(18, 'https://fabrika-stil.ru/_upload/models/site/big/PN-2/PN-21631614303.jpg', true, NOW() - INTERVAL '22 days'),

(19, 'https://www.unitex.ru/image_cache/517x305_sized_-image_products-chairs-personal-210250.jpg', true, NOW() - INTERVAL '22 days'),
(20, 'https://m.pm.ru/global_images/goods/274/19c/19c/02c/1325436_original.jpg', true, NOW() - INTERVAL '22 days'),
(21, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTf_Ew9KOO9nZWltfow5mc8FQ5ItsflUR_04Q&s', true, NOW() - INTERVAL '22 days'),

(22, 'https://avatars.mds.yandex.net/get-mpic/14622206/2a00000196a4a59a63a172fc190875a28b3c/orig', true, NOW() - INTERVAL '22 days'),
(23, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRGekVEexhgAfyCfVMXsfSYspWHJP_Cr31XeQ&s', true, NOW() - INTERVAL '22 days'),
(24, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQeWxCun4cVWW64RSG1I6jyx9pyCT2L0VCrrg&s', true, NOW() - INTERVAL '22 days'),
(25, 'https://gardendecor.ru/pictures/category/small/43407.jpg', true, NOW() - INTERVAL '22 days'),

(26, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTSmDYq356LiUfGhZf2JLBS_GxTz15FeomGXQ&s', true, NOW() - INTERVAL '22 days'),
(27, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQ4fjGS24uAWFIkr1tJewvPPI14utyz4zysKg&s', true, NOW() - INTERVAL '22 days'),
(28, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQdxmxRA3lxdomMPkG0-MmIyuURHXKV-wvQrA&s', true, NOW() - INTERVAL '22 days'),

(29, 'https://cdn.express-office.ru/3194512-900x600/divan-3-kh-mestnyj-marko-ekokozha-ekoteks-3001-chernaya.jpg', true, NOW() - INTERVAL '22 days'),
(30, 'https://basket-05.wbbasket.ru/vol813/part81333/81333529/images/big/1.webp', true, NOW() - INTERVAL '22 days'),
(31, 'https://uno-mebel.ru/image/cache/catalog/products/Naturalnaya-sosna/shkaf-dlya-posudy-far-125%283%29-auto_width_1000.jpg', true, NOW() - INTERVAL '22 days'),
(32, 'https://s4.lm-cdn.ru/photo/mebel_dlya_prihozhey/veshalki_napolnye/napolnaya_veshalka_piko_4_seryy_beton/napolnaya_veshalka_piko_4_seryy_beton_276_414_1_0_5667267.jpg', true, NOW() - INTERVAL '22 days'),
(33, 'https://mypuff.ru/u/2_1730355213_1730355213_1.jpg', true, NOW() - INTERVAL '22 days'),
(34, 'https://taipit-mebel.ru/upload/resize_cache/iblock/e43/xohqqgp60d2c4j1nn6hz8091qo0u29yc/1175_710_1/igrovoe-kompyuternoe-kreslo-chairman-ch26-tkan-seryy.jpg', true, NOW() - INTERVAL '22 days'),
(35, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTI2EwdjdpCTwFtW1Iw6RAn0ZrJJfskjwo3Xw&s', true, NOW() - INTERVAL '22 days'),
(36, 'https://mkub.ru/image/cache/catalog/Malych/d314bf6d-b1b9-4693-ba51-91507c5d3ffb-1024x768.jpg', true, NOW() - INTERVAL '22 days');

-- Добавление записей в таблицу избранного
INSERT INTO favorites (user_id, listing_id, created_at) VALUES
-- Избранное пользователя 1
(1, 5, NOW() - INTERVAL '27 days'),
(1, 8, NOW() - INTERVAL '26 days'),
(1, 15, NOW() - INTERVAL '25 days'),

-- Избранное пользователя 2
(2, 1, NOW() - INTERVAL '26 days'),
(2, 12, NOW() - INTERVAL '25 days'),
(2, 18, NOW() - INTERVAL '24 days'),

-- Избранное пользователя 3
(3, 4, NOW() - INTERVAL '25 days'),
(3, 7, NOW() - INTERVAL '24 days'),
(3, 16, NOW() - INTERVAL '23 days'),

-- Избранное пользователя 4
(4, 2, NOW() - INTERVAL '24 days'),
(4, 9, NOW() - INTERVAL '23 days'),
(4, 23, NOW() - INTERVAL '22 days'),

-- Избранное пользователя 5
(5, 1, NOW() - INTERVAL '23 days'),
(5, 12, NOW() - INTERVAL '22 days'),
(5, 22, NOW() - INTERVAL '21 days'),

-- Избранное пользователя 6
(6, 8, NOW() - INTERVAL '22 days'),
(6, 15, NOW() - INTERVAL '21 days'),
(6, 20, NOW() - INTERVAL '20 days'),

-- Избранное пользователя 7
(7, 2, NOW() - INTERVAL '21 days'),
(7, 13, NOW() - INTERVAL '20 days'),
(7, 18, NOW() - INTERVAL '19 days'),

-- Избранное пользователя 8
(8, 4, NOW() - INTERVAL '20 days'),
(8, 11, NOW() - INTERVAL '19 days'),
(8, 15, NOW() - INTERVAL '18 days'),

-- Избранное пользователя 9
(9, 6, NOW() - INTERVAL '19 days'),
(9, 14, NOW() - INTERVAL '18 days'),
(9, 20, NOW() - INTERVAL '17 days'),

-- Избранное пользователя 10
(10, 3, NOW() - INTERVAL '18 days'),
(10, 9, NOW() - INTERVAL '17 days'),
(10, 16, NOW() - INTERVAL '16 days');

-- Добавление записей о покупках
INSERT INTO purchases (listing_id, buyer_id, seller_id, price, purchased_at) VALUES
-- Покупки, где продавец - пользователь 1
(3, 4, 1, 18000, NOW() - INTERVAL '20 days'),

-- Покупки, где продавец - пользователь 2
(6, 3, 2, 6000, NOW() - INTERVAL '21 days'),

-- Покупки, где продавец - пользователь 3
(10, 5, 3, 22000, NOW() - INTERVAL '19 days'),

-- Покупки, где продавец - пользователь 4
(14, 1, 4, 7000, NOW() - INTERVAL '18 days'),

-- Покупки, где продавец - пользователь 5
(17, 6, 5, 6500, NOW() - INTERVAL '17 days'),

-- Покупки, где продавец - пользователь 6
(20, 7, 6, 11000, NOW() - INTERVAL '16 days'),

-- Покупки, где продавец - пользователь 7
(22, 8, 7, 9500, NOW() - INTERVAL '15 days'),

-- Покупки, где продавец - пользователь 8
(25, 9, 8, 17000, NOW() - INTERVAL '14 days'),

-- Покупки, где продавец - пользователь 9
(28, 10, 9, 9800, NOW() - INTERVAL '13 days'),

-- Покупки, где продавец - пользователь 10
(30, 2, 10, 6500, NOW() - INTERVAL '12 days');

-- Обновление статуса проданных объявлений
UPDATE listings SET status = 'sold' WHERE id IN (3, 6, 10, 14, 17, 20, 22, 25, 28, 30);

-- Создание чатов между пользователями
INSERT INTO chats (listing_id, buyer_id, seller_id, created_at) VALUES
-- Чаты, связанные с объявлениями
(1, 2, 1, NOW() - INTERVAL '28 days'),
(5, 1, 2, NOW() - INTERVAL '27 days'),
(8, 1, 3, NOW() - INTERVAL '26 days'),
(12, 2, 4, NOW() - INTERVAL '25 days'),
(15, 3, 5, NOW() - INTERVAL '24 days'),
(18, 4, 6, NOW() - INTERVAL '23 days'),
(21, 5, 7, NOW() - INTERVAL '22 days'),
(24, 6, 8, NOW() - INTERVAL '21 days'),
(27, 7, 9, NOW() - INTERVAL '20 days'),
(31, 8, 10, NOW() - INTERVAL '19 days'),

-- Чаты для проданных объявлений
(3, 4, 1, NOW() - INTERVAL '22 days'),
(6, 3, 2, NOW() - INTERVAL '23 days'),
(10, 5, 3, NOW() - INTERVAL '21 days'),
(14, 1, 4, NOW() - INTERVAL '20 days'),
(17, 6, 5, NOW() - INTERVAL '19 days');

-- Добавление сообщений в чаты
INSERT INTO messages (chat_id, user_id, content, created_at, is_read) VALUES
-- Сообщения в чате между пользователями 2 и 1 (объявление 1)
(1, 2, 'Здравствуйте, диван еще доступен?', NOW() - INTERVAL '28 days', true),
(1, 1, 'Да, конечно! Вы можете приехать посмотреть.', NOW() - INTERVAL '28 days', true),
(1, 2, 'Какие у него размеры?', NOW() - INTERVAL '27 days', true),
(1, 1, 'Длина 2.2м, ширина 1м. Очень удобный.', NOW() - INTERVAL '27 days', true),
(1, 2, 'Спасибо, подумаю и свяжусь с вами.', NOW() - INTERVAL '27 days', true),

-- Сообщения в чате между пользователями 1 и 2 (объявление 5)
(2, 1, 'Добрый день, интересует ваше кресло-качалка.', NOW() - INTERVAL '27 days', true),
(2, 2, 'Здравствуйте! Да, оно в отличном состоянии.', NOW() - INTERVAL '27 days', true),
(2, 1, 'Можно ли договориться о скидке?', NOW() - INTERVAL '26 days', true),
(2, 2, 'Могу скинуть 500 рублей, больше не могу.', NOW() - INTERVAL '26 days', true),
(2, 1, 'Хорошо, меня устраивает. Когда можно забрать?', NOW() - INTERVAL '26 days', false),

-- Сообщения в чате между пользователями 1 и 3 (объявление 8)
(3, 1, 'Привет, диван-кровать в хорошем состоянии?', NOW() - INTERVAL '26 days', true),
(3, 3, 'Привет! Да, механизм работает отлично.', NOW() - INTERVAL '26 days', true),
(3, 1, 'Сколько лет им пользовались?', NOW() - INTERVAL '25 days', true),
(3, 3, 'Около 3 лет, но очень аккуратно.', NOW() - INTERVAL '25 days', false),

-- Сообщения в чате между пользователями 4 и 1 (объявление 3 - проданное)
(11, 4, 'Здравствуйте, могу ли я посмотреть кровать сегодня?', NOW() - INTERVAL '22 days', true),
(11, 1, 'Да, конечно. Можете подъехать после 18:00', NOW() - INTERVAL '22 days', true),
(11, 4, 'Отлично, буду примерно в 19:00', NOW() - INTERVAL '21 days', true),
(11, 1, 'Хорошо, жду. Адрес скину в личку.', NOW() - INTERVAL '21 days', true),
(11, 4, 'Кровать отличная, беру! Завтра переведу деньги.', NOW() - INTERVAL '20 days', true),
(11, 1, 'Договорились! Спасибо за покупку.', NOW() - INTERVAL '20 days', true);

-- Обновляем is_read для некоторых сообщений, чтобы у пользователей были непрочитанные сообщения
UPDATE messages SET is_read = false WHERE chat_id = 2 AND user_id = 1 AND created_at > NOW() - INTERVAL '26 days';
UPDATE messages SET is_read = false WHERE chat_id = 3 AND user_id = 3 AND created_at > NOW() - INTERVAL '25 days';

-- Добавление 2ФА кодов для некоторых пользователей (для демонстрационных целей)
INSERT INTO two_factor_codes (user_id, code, expires_at) VALUES
(1, '123456', NOW() + INTERVAL '15 minutes'),
(2, '654321', NOW() + INTERVAL '15 minutes'),
(3, '789012', NOW() + INTERVAL '15 minutes'),
(4, '456789', NOW() + INTERVAL '15 minutes'),
(5, '321654', NOW() + INTERVAL '15 minutes');

-- Создание дополнительных чатов между пользователями (без привязки к объявлениям)
INSERT INTO chats (buyer_id, seller_id, created_at) VALUES
(1, 3, NOW() - INTERVAL '24 days'),
(2, 4, NOW() - INTERVAL '23 days'),
(3, 5, NOW() - INTERVAL '22 days'),
(4, 6, NOW() - INTERVAL '21 days'),
(5, 7, NOW() - INTERVAL '20 days');

-- Добавление сообщений в дополнительные чаты
INSERT INTO messages (chat_id, user_id, content, created_at, is_read) VALUES
-- Сообщения в чате между пользователями 1 и 3 (без объявления)
(16, 1, 'Привет! Видел твои объявления, интересные предложения.', NOW() - INTERVAL '24 days', true),
(16, 3, 'Спасибо! Что-то приглянулось конкретно?', NOW() - INTERVAL '24 days', true),
(16, 1, 'Пока смотрю, выбираю. Как долго ты занимаешься продажей мебели?', NOW() - INTERVAL '23 days', true),
(16, 3, 'Около года. Всю мебель проверяю лично перед продажей.', NOW() - INTERVAL '23 days', true),

-- Сообщения в чате между пользователями 2 и 4
(17, 2, 'Здравствуйте! Вы продаёте детскую кроватку?', NOW() - INTERVAL '23 days', true),
(17, 4, 'Да, она всё ещё доступна.', NOW() - INTERVAL '23 days', true),
(17, 2, 'Какие у неё размеры?', NOW() - INTERVAL '22 days', true),
(17, 4, 'Стандартные - 60x120 см. Матрас в комплекте.', NOW() - INTERVAL '22 days', false),

-- Сообщения в чате между пользователями 3 и 5
(18, 3, 'Добрый день! Интересует полка настенная.', NOW() - INTERVAL '22 days', true),
(18, 5, 'Добрый! Да, полка в наличии.', NOW() - INTERVAL '22 days', true),
(18, 3, 'Какой максимальный вес она выдерживает?', NOW() - INTERVAL '21 days', true),
(18, 5, 'До 15 кг. Крепления в комплекте.', NOW() - INTERVAL '21 days', true);

-- Добавим больше избранных объявлений для разнообразия
INSERT INTO favorites (user_id, listing_id, created_at) VALUES
(1, 24, NOW() - INTERVAL '19 days'),
(1, 31, NOW() - INTERVAL '18 days'),
(2, 9, NOW() - INTERVAL '19 days'),
(2, 21, NOW() - INTERVAL '18 days'),
(3, 13, NOW() - INTERVAL '19 days'),
(4, 18, NOW() - INTERVAL '18 days'),
(5, 29, NOW() - INTERVAL '17 days'),
(6, 31, NOW() - INTERVAL '16 days'),
(7, 27, NOW() - INTERVAL '15 days'),
(8, 18, NOW() - INTERVAL '14 days'),
(9, 4, NOW() - INTERVAL '13 days'),
(10, 7, NOW() - INTERVAL '12 days');

-- Добавление сообщений о завершении покупок
INSERT INTO messages (chat_id, user_id, content, created_at, is_read) VALUES
-- Сообщения для последних покупок
(1, 2, 'Здравствуйте! Я решил купить диван, когда можно забрать?', NOW() - INTERVAL '11 days', true),
(1, 1, 'Добрый день! В любое время, я дома весь день.', NOW() - INTERVAL '11 days', true),
(1, 2, 'Отлично, буду у вас завтра в 16:00', NOW() - INTERVAL '10 days', true),
(1, 1, 'Хорошо, буду ждать. Спасибо за покупку!', NOW() - INTERVAL '10 days', true);

-- Статистические запросы для проверки данных
-- SELECT 'Пользователи:', COUNT(*) FROM users;
-- SELECT 'Объявления:', COUNT(*) FROM listings;
-- SELECT 'Изображения:', COUNT(*) FROM listing_images;
-- SELECT 'Избранное:', COUNT(*) FROM favorites;
-- SELECT 'Чаты:', COUNT(*) FROM chats;
-- SELECT 'Сообщения:', COUNT(*) FROM messages;
-- SELECT 'Покупки:', COUNT(*) FROM purchases;
-- SELECT 'Категории:', COUNT(*) FROM categories;
-- SELECT 'Коды 2ФА:', COUNT(*) FROM two_factor_codes;

-- Выводим примеры данных для проверки
-- SELECT 'Пример пользователя:', * FROM users LIMIT 1;
-- SELECT 'Пример объявления:', * FROM listings LIMIT 1;
-- SELECT 'Пример изображения:', * FROM listing_images LIMIT 1;
-- SELECT 'Пример избранного:', * FROM favorites LIMIT 1;
-- SELECT 'Пример чата:', * FROM chats LIMIT 1;
-- SELECT 'Пример сообщения:', * FROM messages LIMIT 1;
-- SELECT 'Пример покупки:', * FROM purchases LIMIT 1;