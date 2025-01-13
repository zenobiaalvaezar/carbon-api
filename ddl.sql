CREATE TABLE roles (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL
);

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  role_id INT NOT NULL,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  phone VARCHAR(255),
  address TEXT,
  created_at DATE NOT NULL,
  FOREIGN KEY (role_id) REFERENCES roles (id)
);

CREATE TABLE fuels (
  id SERIAL PRIMARY KEY,
  category VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  emission_factor FLOAT NOT NULL,
  price FLOAT NOT NULL,
  unit VARCHAR(255) NOT NULL
);

CREATE TABLE carbon_fuels (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  fuel_id INT NOT NULL,
  price FLOAT NOT NULL,
  usage_amount FLOAT NOT NULL,
  usage_type VARCHAR(255) NOT NULL,
  total_consumption FLOAT NOT NULL,
  emission_factor FLOAT NOT NULL,
  emission_amount FLOAT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users (id),
  FOREIGN KEY (fuel_id) REFERENCES fuels (id)
);

CREATE TABLE electrics (
  id SERIAL PRIMARY KEY,
  province VARCHAR(255) NOT NULL,
  emission_factor FLOAT NOT NULL,
  price FLOAT NOT NULL
);

CREATE TABLE carbon_electrics (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  electric_id INT NOT NULL,
  price FLOAT NOT NULL,
  usage_amount FLOAT NOT NULL,
  usage_type VARCHAR(255) NOT NULL,
  total_consumption FLOAT NOT NULL,
  emission_factor FLOAT NOT NULL,
  emission_amount FLOAT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users (id),
  FOREIGN KEY (electric_id) REFERENCES electrics (id)
);

CREATE TABLE carbon_summaries (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  fuel_emission FLOAT NOT NULL,
  electric_emission FLOAT NOT NULL,
  total_emission FLOAT NOT NULL,
  total_tree INT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE tree_categories (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL
);

CREATE TABLE trees (
  id SERIAL PRIMARY KEY,
  tree_category_id INT NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  price FLOAT NOT NULL,
  stock INT NOT NULL,
  FOREIGN KEY (tree_category_id) REFERENCES tree_categories (id)
);

CREATE TABLE carts (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  tree_id INT NOT NULL,
  quantity INT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users (id),
  FOREIGN KEY (tree_id) REFERENCES trees (id)
);

CREATE TABLE transactions (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  total_price FLOAT NOT NULL,
  created_at DATE NOT NULL,
  payment_method VARCHAR(255) NOT NULL,
  payment_status VARCHAR(255) NOT NULL,
  payment_at DATE,
  FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE transaction_details (
  id SERIAL PRIMARY KEY,
  transaction_id INT NOT NULL,
  tree_id INT NOT NULL,
  quantity INT NOT NULL,
  price FLOAT NOT NULL,
  total_price FLOAT NOT NULL,
  FOREIGN KEY (transaction_id) REFERENCES transactions (id),
  FOREIGN KEY (tree_id) REFERENCES trees (id)
);

INSERT INTO roles (name)
VALUES
('admin'),
('customer');

-- insert users with password "12345"
INSERT INTO users (role_id, name, email, password, phone, address, created_at)
VALUES
((SELECT id FROM roles WHERE name='admin'), 'admin', 'admin@gmail.com', '$2a$10$8BBI3BLnTfrvwg6WiISSPeUf.NGpwDhEfKPCTHFfDbD/uLw3lI8BK', '08123456789', 'Jakarta', '2021-01-01'),
((SELECT id FROM roles WHERE name='customer'), 'customer', 'customer@gmail.com', '$2a$10$8BBI3BLnTfrvwg6WiISSPeUf.NGpwDhEfKPCTHFfDbD/uLw3lI8BK', '08123456789', 'Bekasi', '2021-01-02');

INSERT INTO fuels (category, name, emission_factor, price, unit)
VALUES
('Bahan Bakar Cair', 'Pertamax Plus/Turbo', 2.36, 13250, 'Liter'),
('Bahan Bakar Cair', 'Pertamax', 2.36, 12500, 'Liter'),
('Bahan Bakar Cair', 'Pertalite', 2.36, 12000, 'Liter'),
('Bahan Bakar Cair', 'Premium', 2.37, 6500, 'Liter'),
('Bahan Bakar Cair', 'Avtur', 2.61, 17272, 'Liter'),
('Bahan Bakar Cair', 'Minyak Tanah', 2.61, 20000, 'Liter'),
('Bahan Bakar Cair', 'Pertadex', 2.69, 9000, 'Liter'),
('Bahan Bakar Cair', 'Dexlite', 2.71, 13500, 'Liter'),
('Bahan Bakar Cair', 'Bio Solar', 2.74, 12100, 'Liter'),
('Bahan Bakar Cair', 'Minyak Diesel (MDF)', 2.84, 10000, 'Liter'),
('Bahan Bakar Cair', 'Minyak Bakar (MFO)', 3.17, 10000, 'Liter'),
('Bahan Bakar Padat', 'Batubara', 2.02, 13211, 'kg'),
('Bahan Bakar Padat', 'Briket Batubara', 2.06, 16950, 'kg'),
('Bahan Bakar Padat', 'Arang', 3.38, 13150, 'kg'),
('Bahan Bakar Gas', 'Gas Alam', 2.20, 13200, 'NM3'),
('Bahan Bakar Gas', 'LPG', 3.09, 6800, 'kg'),
('Bahan Bakar Gas', 'LGV', 3.07, 11000, 'kg'),
('Bahan Bakar Gas', 'LNG', 2.76, 9000, 'kg');

INSERT INTO electrics (province, emission_factor, price)
VALUES
('Nanggroe Aceh Darussalam', 0.79, 1450),
('Sumatera Utara', 0.79, 1450),
('Sumatera Selatan', 0.79, 1450),
('Sumatera Barat', 0.72, 1450),
('Bengkulu', 0.70, 1450),
('Riau', 0.60, 1450),
('Kepulauan Riau', 0.61, 1450),
('Jambi', 0.79, 1450),
('Lampung', 0.79, 1450),
('Bangka Belitung', 0.86, 1450),
('Kalimantan Barat', 1.71, 1450),
('Kalimantan Timur', 1.15, 1450),
('Kalimantan Selatan', 1.23, 1450),
('Kalimantan Tengah', 1.23, 1450),
('Kalimantan Utara', 0.30, 1450),
('Banten', 0.82, 1450),
('DKI Jakarta', 0.82, 1450),
('Jawa Barat', 0.82, 1450),
('Jawa Tengah', 0.73, 1450),
('Di Yogyakarta', 0.82, 1450),
('Jawa Timur', 0.82, 1450),
('Bali', 0.53, 1450),
('Nusa Tenggara Timur', 0.60, 1450),
('Nusa Tenggara Barat', 0.86, 1450),
('Gorontalo', 0.62, 1450),
('Sulawesi Barat', 0.79, 1450),
('Sulawesi Tengah', 0.62, 1450),
('Sulawesi Utara', 0.73, 1450),
('Sulawesi Tenggara', 0.66, 1450),
('Sulawesi Selatan', 0.79, 1450),
('Maluku Utara', 0.63, 1450),
('Maluku', 0.66, 1450),
('Papua Barat', 0.57, 1450),
('Papua', 0.57, 1450),
('Papua Tengah', 0.57, 1450),
('Papua Pegunungan', 0.57, 1450),
('Papua Selatan', 0.57, 1450),
('Papua Barat Daya', 0.57, 1450);

INSERT INTO tree_categories (name)
VALUES
('Pohon Kayu Keras'),
('Pohon Peneduh & Pengikat Air'),
('Pohon Pelindung Satwa Liar'),
('Pohon Cepat Tumbuh'),
('Pohon Endemik Lokal');

INSERT INTO trees (tree_category_id, name, description, price, stock)
VALUES
((SELECT id FROM tree_categories WHERE name='Pohon Kayu Keras'), 'Jati', 'Kayu berkualitas tinggi, tahan lama, sering digunakan untuk furnitur mewah.', 50000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Kayu Keras'), 'Mahoni', 'Kayu halus, mudah diolah, cocok untuk mebel dan kayu lapis.', 30000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Kayu Keras'), 'Meranti', 'Tahan terhadap cuaca, banyak digunakan dalam konstruksi dan dekorasi rumah.', 40000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Kayu Keras'), 'Ulin', 'Kayu super kuat, tahan air, sering digunakan untuk bahan bangunan berat.', 100000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Kayu Keras'), 'Sengon Buto', 'Cepat tumbuh, cocok untuk penahan angin dan pelindung tanah dari erosi.', 25000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Peneduh & Pengikat Air'), 'Bambu', 'Cepat tumbuh, kuat, dan efektif sebagai pengikat air di tanah yang miring.', 15000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Peneduh & Pengikat Air'), 'Kepuh', 'Akar kuat untuk mencegah longsor dan memberikan naungan yang lebat.', 35000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Peneduh & Pengikat Air'), 'Beringin', 'Ikonik sebagai pohon pelindung dengan akar gantung yang unik.', 50000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Peneduh & Pengikat Air'), 'Angsana', 'Cepat tumbuh, memberikan naungan, dan memiliki bunga berwarna kuning cerah.', 30000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Peneduh & Pengikat Air'), 'Glirisidia', 'Sumber pupuk hijau alami, cocok untuk agroforestri dan konservasi lahan.', 20000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Pelindung Satwa Liar'), 'Kemiri', 'Menghasilkan biji berharga dan tempat tinggal bagi burung.', 25000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Pelindung Satwa Liar'), 'Durian', 'Menghasilkan buah, memberikan makanan untuk satwa liar dan manusia.', 70000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Pelindung Satwa Liar'), 'Karet', 'Menghasilkan getah karet yang bernilai ekonomi tinggi.', 40000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Pelindung Satwa Liar'), 'Ketapang', 'Daunnya bisa menurunkan pH air, ideal untuk habitat ikan air tawar.', 35000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Pelindung Satwa Liar'), 'Trembesi', 'Kanopi lebar, efektif menyerap karbon dioksida, cocok untuk penghijauan kota.', 45000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Cepat Tumbuh'), 'Sengon', 'Kayu ringan, cocok untuk produksi papan dan pulp.', 20000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Cepat Tumbuh'), 'Akasia', 'Cepat tumbuh, digunakan untuk kayu bakar, pulp, dan reboisasi lahan kritis.', 25000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Cepat Tumbuh'), 'Eukaliptus', 'Kayu serbaguna dan menghasilkan minyak atsiri bernilai ekonomi.', 30000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Cepat Tumbuh'), 'Jabon', 'Kayu putih halus, ideal untuk industri kayu ringan seperti furnitur.', 30000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Cepat Tumbuh'), 'Petai Cina', 'Penyubur tanah alami, menghasilkan biji untuk konsumsi.', 15000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Endemik Lokal'), 'Rasamala', 'Kayu wangi dan kuat, digunakan untuk bahan konstruksi tradisional.', 50000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Endemik Lokal'), 'Pulai', 'Akar kuat, cocok untuk penahan tanah dan konservasi lahan.', 40000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Endemik Lokal'), 'Damar', 'Menghasilkan resin damar bernilai tinggi untuk industri.', 60000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Endemik Lokal'), 'Keruing', 'Kayu berkualitas tinggi untuk bahan bangunan dan furnitur.', 70000, 5),
((SELECT id FROM tree_categories WHERE name='Pohon Endemik Lokal'), 'Merbau', 'Kayu keras yang tahan lama, sering digunakan untuk lantai dan furnitur.', 80000, 5);
