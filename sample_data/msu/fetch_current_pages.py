#!/usr/bin/env python3
"""
Temporary script to fetch current MSU faculty pages for testing quota parsing.
This script fetches a few key pages to analyze quota application parsing issues.
"""

import requests
import time
import os
from urllib.parse import urlparse

# Key URLs to fetch for quota analysis
URLS_TO_FETCH = [
    "https://cpk.msu.ru/rating/dep_11",  # История, История искусств, История международных отношений
    "https://cpk.msu.ru/rating/dep_31",  # Политология, Управление персоналом 
    "https://cpk.msu.ru/rating/dep_01",  # Математика, Механика
    "https://cpk.msu.ru/rating/dep_08",  # Науки о земле, Туризм
    "https://cpk.msu.ru/rating/dep_05",  # Фундаментальная и прикладная биология
]

def fetch_and_save_page(url: str, output_dir: str):
    """Fetch a page and save it as HTML file."""
    try:
        print(f"Fetching {url}...")
        response = requests.get(url, timeout=30)
        response.raise_for_status()
        
        # Extract filename from URL
        parsed = urlparse(url)
        filename = parsed.path.split('/')[-1] + '.html'
        filepath = os.path.join(output_dir, filename)
        
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(response.text)
        
        print(f"Saved {filename} ({len(response.text)} characters)")
        return True
        
    except Exception as e:
        print(f"Error fetching {url}: {e}")
        return False

def main():
    script_dir = os.path.dirname(os.path.abspath(__file__))
    
    print("Fetching current MSU faculty pages for quota analysis...")
    print(f"Output directory: {script_dir}")
    
    success_count = 0
    
    for url in URLS_TO_FETCH:
        if fetch_and_save_page(url, script_dir):
            success_count += 1
        
        # Be polite to the server
        time.sleep(2)
    
    print(f"\nCompleted: {success_count}/{len(URLS_TO_FETCH)} pages fetched successfully")

if __name__ == "__main__":
    main()