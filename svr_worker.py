# svr_worker.py
import sys
import json
import pandas as pd
from datetime import datetime, timedelta
from sklearn.svm import SVR
from sklearn.preprocessing import StandardScaler

def main():
    if len(sys.argv) < 2:
        print(json.dumps({"error": "No arguments"}))
        return
    
    input_user = json.loads(sys.argv[1]) 
    
    jawaban_rescaled = [(x - 1) * (3 / 4) for x in input_user]
    
    try:
        df = pd.read_csv('PHQ9_GAD7_df_comma.csv', sep=',')
    except FileNotFoundError:
        print(json.dumps({"error": "Dataset CSV tidak ditemukan"}))
        return

   
    if 'Tanggal' not in df.columns:
        import numpy as np
        tgl_sekarang = datetime.now()
        df['Tanggal'] = [(tgl_sekarang - timedelta(days=int(np.random.randint(0, 40)))).strftime('%Y-%m-%d') for _ in range(len(df))]

    df['Tanggal'] = pd.to_datetime(df['Tanggal'])
    batas_satu_bulan = datetime.now() - timedelta(days=30)
    
    df_filtered = df[df['Tanggal'] >= batas_satu_bulan]
    
    if len(df_filtered) < 3:
        print(json.dumps({"error": "Data satu bulan terakhir kurang dari batas minimum"}))
        return

    X = df_filtered[['PHQ1', 'PHQ2', 'GAD1', 'GAD2']]
    y = df_filtered['GAD7'] 
    
    scaler = StandardScaler()
    X_scaled = scaler.fit_transform(X)
    
    model_svr = SVR(kernel='rbf', C=1.5, epsilon=0.1)
    model_svr.fit(X_scaled, y)
    
    X_user_scaled = scaler.transform([jawaban_rescaled])
    prediksi_target = model_svr.predict(X_user_scaled)[0]
    
    prediksi_cli = (prediksi_target / 3) * (5 - 1) + 1
    prediksi_cli = max(1.0, min(5.0, prediksi_cli))

    output = {
        "prediksi_skor_q5": round(prediksi_cli, 2),
        "status": "success",
        "data_points_used": len(df_filtered)
    }
    print(json.dumps(output))

if __name__ == '__main__':
    main()