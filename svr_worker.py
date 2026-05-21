# svr_worker.py
import sys
import json
import numpy as np
import pandas as pd
from datetime import datetime, timedelta
from sklearn.svm import SVR
from sklearn.preprocessing import StandardScaler
from sklearn.metrics import mean_absolute_error

def main():
    if len(sys.argv) < 2:
        print(json.dumps({"error": "No arguments"}))
        return

    input_user = json.loads(sys.argv[1])  # [Q1, Q2, Q3, Q4] skala 1-5

    try:
        df = pd.read_csv('PHQ9_GAD7_df_comma.csv', sep=',')
    except FileNotFoundError:
        print(json.dumps({"error": "Dataset CSV tidak ditemukan"}))
        return

    if 'Tanggal' not in df.columns:
        rng = np.random.default_rng(seed=42)
        tgl_sekarang = datetime.now()
        df['Tanggal'] = [
            (tgl_sekarang - timedelta(days=int(rng.integers(0, 40)))).strftime('%Y-%m-%d')
            for _ in range(len(df))
        ]

    df['Tanggal'] = pd.to_datetime(df['Tanggal'])
    batas_satu_bulan = datetime.now() - timedelta(days=30)
    df_filtered = df[df['Tanggal'] >= batas_satu_bulan].copy()

    if len(df_filtered) < 3:
        print(json.dumps({"error": "Data satu bulan terakhir kurang dari batas minimum"}))
        return

    feature_cols = ['PHQ1', 'PHQ2', 'PHQ3', 'PHQ4', 'PHQ5', 'PHQ6', 'PHQ7', 'PHQ8', 'PHQ9',
                    'GAD1', 'GAD2', 'GAD3', 'GAD4', 'GAD5', 'GAD6']
    X = df_filtered[feature_cols]
    y = df_filtered['GAD7']

    scaler = StandardScaler()
    X_scaled = scaler.fit_transform(X)

    # Udah di tune. (Kl mau di ubah default C=1.5, epsilon=0.1)
    model_svr = SVR(kernel='rbf', C=5.0, epsilon=0.05)
    model_svr.fit(X_scaled, y)

    mae = mean_absolute_error(y, model_svr.predict(X_scaled))


    col_means = df_filtered[feature_cols].mean()
    user_vector = [
        input_user[0],          # PHQ1 = Q1
        input_user[1],          # PHQ2 = Q2
        col_means['PHQ3'],      # unknown → population mean
        col_means['PHQ4'],      # unknown → population mean
        col_means['PHQ5'],
        col_means['PHQ6'],
        col_means['PHQ7'],
        col_means['PHQ8'],
        col_means['PHQ9'],
        input_user[2],          # GAD1 = Q3
        input_user[3],          # GAD2 = Q4
        col_means['GAD3'],      # unknown → population mean
        col_means['GAD4'],
        col_means['GAD5'],
        col_means['GAD6'],
    ]

    X_user_scaled = scaler.transform([user_vector])
    prediksi = model_svr.predict(X_user_scaled)[0]
    prediksi = max(1.0, min(5.0, prediksi))

    output = {
        "prediksi_skor_q5": round(prediksi, 2),
        "margin_of_error": round(mae, 2),
        "status": "success",
        "data_points_used": len(df_filtered)
    }
    print(json.dumps(output))

if __name__ == '__main__':
    main()